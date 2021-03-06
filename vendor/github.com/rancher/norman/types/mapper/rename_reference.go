package mapper

import (
	"strings"

	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/definition"
)

type RenameReference struct {
	mapper types.Mapper
}

func (r *RenameReference) FromInternal(data map[string]interface{}) {
	if r.mapper != nil {
		r.mapper.FromInternal(data)
	}
}

func (r *RenameReference) ToInternal(data map[string]interface{}) {
	if r.mapper != nil {
		r.mapper.ToInternal(data)
	}
}

func (r *RenameReference) ModifySchema(schema *types.Schema, schemas *types.Schemas) error {
	var mappers []types.Mapper
	for name, field := range schema.ResourceFields {
		if definition.IsReferenceType(field.Type) && strings.HasSuffix(name, "Name") {
			newName := strings.TrimSuffix(name, "Name") + "Id"
			move := Move{From: name, To: newName}
			if err := move.ModifySchema(schema, schemas); err != nil {
				return err
			}

			mappers = append(mappers, move)
		} else if definition.IsArrayType(field.Type) && definition.IsReferenceType(definition.SubType(field.Type)) && strings.HasSuffix(name, "Names") {
			newName := strings.TrimSuffix(name, "Names") + "Ids"
			move := Move{From: name, To: newName}
			if err := move.ModifySchema(schema, schemas); err != nil {
				return err
			}

			mappers = append(mappers, move)
		}
	}

	if len(mappers) > 0 {
		r.mapper = types.Mappers(mappers)
	}

	return nil
}
