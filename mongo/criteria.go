package mongo

import (
	"errors"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/go-specification"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

var ComparisonConversion = map[specification.Comparator]string{
	specification.ComparisonEq:  "$eq",
	specification.ComparisonGt:  "$gt",
	specification.ComparisonGte: "$gte",
	specification.ComparisonLt:  "$lt",
	specification.ComparisonLte: "$lte",
	specification.ComparisonNe:  "$ne",
}

type Criteria interface {
	GetExpression() bson.M
}

type Attr struct {
	Name       string
	Value      any
	Comparison specification.Comparator
}

func (a Attr) GetExpression() bson.M {
	return bson.M{a.Name: bson.M{ComparisonConversion[a.Comparison]: a.Value}}
}

type MongoAnd struct {
	Operands []Criteria
}

func (a MongoAnd) GetExpression() bson.M {
	var expressions []bson.M
	for _, operand := range a.Operands {
		expressions = append(expressions, operand.GetExpression())
	}

	return bson.M{"$and": expressions}
}

type MongoOr struct {
	Operands []Criteria
}

func (o MongoOr) GetExpression() bson.M {
	var expressions []bson.M
	for _, operand := range o.Operands {
		expressions = append(expressions, operand.GetExpression())
	}

	return bson.M{"$or": expressions}
}

type MongoNot struct {
	Operand Criteria
}

func (n MongoNot) GetExpression() bson.M {
	return bson.M{"$not": n.Operand.GetExpression()}
}

func ConvertToMongoCriteria(c specification.Criteria, subject entity.Entity) (Criteria, error) {
	switch c.(type) {
	case specification.Attr:
		ca := c.(specification.Attr)

		field, ok := reflect.TypeOf(subject).Elem().FieldByName(ca.Name)
		if !ok {
			return nil, errors.New("field not found")
		}

		tag := field.Tag.Get("bson")
		if tag == "" {
			return nil, errors.New("field is not bson tagged")
		}
		if strings.Contains(tag, ",") {
			tag = strings.Split(tag, ",")[0]
		}

		return Attr{Name: tag, Value: ca.Value, Comparison: ca.Comparison}, nil
	case specification.And:
		var ops []Criteria
		for _, operand := range c.(specification.And).Operands {
			mc, err := ConvertToMongoCriteria(operand, subject)
			if err != nil {
				return nil, err
			}
			ops = append(ops, mc)
		}
		return MongoAnd{Operands: ops}, nil
	case specification.Or:
		var ops []Criteria
		for _, operand := range c.(specification.Or).Operands {
			mc, err := ConvertToMongoCriteria(operand, subject)
			if err != nil {
				return nil, err
			}
			ops = append(ops, mc)
		}
		return MongoOr{Operands: ops}, nil
	case specification.Not:
		mc, err := ConvertToMongoCriteria(c.(specification.Not).Operand, subject)
		if err != nil {
			return nil, err
		}

		return MongoNot{Operand: mc}, nil
	}

	return nil, errors.New("unknown criteria type")
}

// TODO how with dark magic can we make this work?
//type MongoResolver struct {
//	impls map[any]any
//}
//
//func NewMongoResolver() *MongoResolver {
//	return &MongoResolver{
//		impls: map[any]any{
//			"Attr": "MongoAttr",
//		},
//	}
//}

//func (r *MongoResolver) Resolve(c Criteria) Criteria {
//	return r.impls[c.GetType()]
//}
