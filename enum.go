package bcs

import (
	"fmt"
	"reflect"
)

var EnumTypes = make(map[reflect.Type][]reflect.Type)

// NOTE: for now it is not thread-safe as it is assumed that all types are registered upon initialization.
func RegisterEnumType[EnumType any](variant any, variants ...any) {
	variants = append([]any{variant}, variants...)

	fmt.Println("XXX", variants)
	enumT := reflect.TypeOf((*EnumType)(nil)).Elem()

	if enumT.Kind() != reflect.Interface {
		panic(fmt.Errorf("RegisterEnumType: enum type %v is not an interface", enumT))
	}

	variantsT := make([]reflect.Type, 0, len(variants))

	for _, v := range variants {
		variantT := reflect.TypeOf(v)

		if variantT.Kind() == reflect.Interface {
			panic(fmt.Errorf("RegisterEnumType: variant type %v of enum %v is an interface", variantT, enumT))
		}

		if !variantT.Implements(enumT) {
			panic(fmt.Errorf("RegisterEnumType: variant type %v does not implement enum %v", variantT, enumT))
		}

		variantsT = append(variantsT, variantT)
	}

	alreadyRegisteredVariants := EnumTypes[enumT]
	if alreadyRegisteredVariants != nil {
		panic(fmt.Errorf("RegisterEnumType: enum type %v is already registered with variants %v", enumT, alreadyRegisteredVariants))
	}

	EnumTypes[enumT] = variantsT
}

func RegisterEnumType1[EnumType any, Variant1 any]() {
	var variant1 Variant1
	RegisterEnumType[EnumType](variant1)
}

func RegisterEnumType2[EnumType any, Variant1 any, Variant2 any]() {
	var variant1 Variant1
	var variant2 Variant2
	RegisterEnumType[EnumType](variant1, variant2)
}

func RegisterEnumType3[EnumType any, Variant1 any, Variant2 any, Variant3 any]() {
	var variant1 Variant1
	var variant2 Variant2
	var variant3 Variant3
	RegisterEnumType[EnumType](variant1, variant2, variant3)
}

func RegisterEnumType4[EnumType any, Variant1 any, Variant2 any, Variant3 any, Variant4 any]() {
	var variant1 Variant1
	var variant2 Variant2
	var variant3 Variant3
	var variant4 Variant4
	RegisterEnumType[EnumType](variant1, variant2, variant3, variant4)
}

func RegisterEnumType5[EnumType any, Variant1 any, Variant2 any, Variant3 any, Variant4 any, Variant5 any]() {
	var variant1 Variant1
	var variant2 Variant2
	var variant3 Variant3
	var variant4 Variant4
	var variant5 Variant5
	RegisterEnumType[EnumType](variant1, variant2, variant3, variant4, variant5)
}

func RegisterEnumType6[EnumType any, Variant1 any, Variant2 any, Variant3 any, Variant4 any, Variant5 any, Variant6 any]() {
	var variant1 Variant1
	var variant2 Variant2
	var variant3 Variant3
	var variant4 Variant4
	var variant5 Variant5
	var variant6 Variant6
	RegisterEnumType[EnumType](variant1, variant2, variant3, variant4, variant5, variant6)
}
