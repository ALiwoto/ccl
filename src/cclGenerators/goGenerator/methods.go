package goGenerator

import (
	"errors"
	"os"
	"strings"

	"github.com/ALiwoto/ccl/src/core/cclErrors"
	"github.com/ALiwoto/ccl/src/core/cclValues"
	"github.com/ALiwoto/ssg/ssg"
)

func (c *GoGenerationContext) GenerateCode() error {
	if c.Options.PackageName == "" {
		return errors.New("package name is required for Go code generation")
	}

	var err error
	err = c.GenerateConstants()
	if err != nil {
		return err
	}

	err = c.GenerateVars()
	if err != nil {
		return err
	}

	err = c.GenerateTypes()
	if err != nil {
		return err
	}

	err = c.GenerateHelpers()
	if err != nil {
		return err
	}

	err = c.GenerateMethods()
	if err != nil {
		return err
	}

	path := c.Options.OutputPath + string(os.PathSeparator)
	if c.ConstantsCode != nil {
		err = ssg.WriteFileStr(path+ConstantsFileName, c.ConstantsCode.String())
		if err != nil {
			return err
		}
	}

	if c.VarsCode != nil {
		err = ssg.WriteFileStr(path+VarsFileName, c.VarsCode.String())
		if err != nil {
			return err
		}
	}

	if c.TypesCode != nil {
		err = ssg.WriteFileStr(path+TypesFileName, c.TypesCode.String())
		if err != nil {
			return err
		}
	}

	if c.HelpersCode != nil {
		err = ssg.WriteFileStr(path+HelpersFileName, c.HelpersCode.String())
		if err != nil {
			return err
		}
	}

	if c.MethodsCode != nil {
		err = ssg.WriteFileStr(path+MethodsFileName, c.MethodsCode.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *GoGenerationContext) GenerateConstants() error {
	c.ConstantsCode = &strings.Builder{}
	c.ConstantsCode.WriteString("// THIS FILE IS AUTOGENERATED BY A CCL TOOL. DO NOT EDIT.\n\n")
	c.ConstantsCode.WriteString("package " + c.Options.PackageName + "\n\n")
	c.ConstantsCode.WriteString("const (\n")

	for _, currentModel := range c.Options.CCLDefinition.Models {
		c.ConstantsCode.WriteString("\tModelId" +
			currentModel.Name + " = " + ssg.ToBase10(currentModel.ModelId) + "\n",
		)
	}

	c.ConstantsCode.WriteString(")\n")
	return nil
}

func (c *GoGenerationContext) GenerateVars() error {
	c.VarsCode = &strings.Builder{}
	c.VarsCode.WriteString("// THIS FILE IS AUTOGENERATED BY A CCL TOOL. DO NOT EDIT.\n\n")
	c.VarsCode.WriteString("package " + c.Options.PackageName + "\n\n")

	return nil
}

func (c *GoGenerationContext) GenerateTypes() error {
	c.TypesCode = &strings.Builder{}
	c.TypesCode.WriteString("// THIS FILE IS AUTOGENERATED BY A CCL TOOL. DO NOT EDIT.\n\n")
	c.TypesCode.WriteString("package " + c.Options.PackageName + "\n\n")

	var neededPackageImports = map[string]bool{}
	for _, currentModel := range c.Options.CCLDefinition.Models {
		for _, currentField := range currentModel.Fields {
			switch currentField.Type {
			case cclValues.TypeNameDateTime:
				if !neededPackageImports["time"] {
					neededPackageImports["time"] = true
				}
			}
		}
	}

	if len(neededPackageImports) > 0 {
		c.TypesCode.WriteString("import (\n")
		for packageName := range neededPackageImports {
			c.TypesCode.WriteString("\t\"" + packageName + "\"\n")
		}
		c.TypesCode.WriteString(")\n\n")
	}

	c.TypesCode.WriteString("type Serializable interface {\n")
	c.TypesCode.WriteString("\tSerializeBinary() ([]byte, error)\n")
	c.TypesCode.WriteString("\tDeserializeBinary(data []byte) error\n")
	c.TypesCode.WriteString("\tCloneEmpty() Serializable\n")
	c.TypesCode.WriteString("}\n\n")

	for _, currentModel := range c.Options.CCLDefinition.Models {
		c.TypesCode.WriteString("type " + currentModel.Name + " struct {\n")
		for _, currentField := range currentModel.Fields {
			theGoType, ok := CCLTypesToGoTypes[currentField.Type]
			customDefinedType := c.Options.CCLDefinition.GetModelByName(currentField.Type)
			if !ok {
				if customDefinedType == nil {
					return &cclErrors.UnsupportedFieldTypeError{
						TypeName:       currentField.Type,
						FieldName:      currentField.Name,
						ModelName:      currentModel.Name,
						TargetLanguage: LanguageName,
					}
				}

				// TODO: add ways to specify this type being pointer or not
				theGoType = "*" + customDefinedType.Name
			}

			// TODO: handle extra operators here
			if currentField.IsArray() {
				theGoType = "[]" + theGoType
			}
			c.TypesCode.WriteString("\t" + currentField.Name + " " + theGoType + "\n")
		}
		c.TypesCode.WriteString("}\n\n")
	}

	return nil
}

func (c *GoGenerationContext) GenerateHelpers() error {
	c.HelpersCode = &strings.Builder{}
	c.HelpersCode.WriteString("// THIS FILE IS AUTOGENERATED BY A CCL TOOL. DO NOT EDIT.\n\n")
	c.HelpersCode.WriteString("package " + c.Options.PackageName + "\n\n")

	return nil
}

func (c *GoGenerationContext) GenerateMethods() error {
	c.MethodsCode = &strings.Builder{}
	c.MethodsCode.WriteString("// THIS FILE IS AUTOGENERATED BY A CCL TOOL. DO NOT EDIT.\n\n")
	c.MethodsCode.WriteString("package " + c.Options.PackageName + "\n\n")

	// default imports
	c.MethodsCode.WriteString("import (\n")
	c.MethodsCode.WriteString("\t\"bytes\"\n")
	c.MethodsCode.WriteString("\t\"encoding/binary\"\n")
	c.MethodsCode.WriteString("\t\"time\"\n")
	c.MethodsCode.WriteString(")\n\n")
	c.MethodsCode.WriteString("const (\n")
	c.MethodsCode.WriteString("\t_ = time.April\n")
	c.MethodsCode.WriteString("\t_ = bytes.MinRead\n")
	c.MethodsCode.WriteString("\t_ = binary.MaxVarintLen16\n")
	c.MethodsCode.WriteString(")\n\n")

	for _, currentModel := range c.Options.CCLDefinition.Models {
		println("currentModel: ", currentModel.Name)
		c.MethodsCode.WriteString("\n//------------------------------------------------------------\n\n")
		c.MethodsCode.WriteString("func (m *" + currentModel.Name + ") GetModelId() int {\n")
		c.MethodsCode.WriteString("\treturn ModelId" + currentModel.Name + "\n")
		c.MethodsCode.WriteString("}\n")

		// generate CloneEmpty() method
		c.MethodsCode.WriteString("\nfunc (m *" + currentModel.Name + ") CloneEmpty() *")
		c.MethodsCode.WriteString(currentModel.Name + " {\n")
		c.MethodsCode.WriteString("\tif m == nil {\n")
		c.MethodsCode.WriteString("\t\treturn nil\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\treturn &" + currentModel.Name + "{}\n")
		c.MethodsCode.WriteString("}\n")

		c.MethodsCode.WriteString("\nfunc (m *" + currentModel.Name + ") CloneEmptySerializable() Serializable {\n")
		c.MethodsCode.WriteString("\tif m == nil {\n")
		c.MethodsCode.WriteString("\t\treturn nil\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\treturn &" + currentModel.Name + "{}\n")
		c.MethodsCode.WriteString("}\n")

		// TODO:
		// This SerializeBinary method generation HAS TO BE MOVED to a attribute handler
		// so users can specify different types of attributes for serialize method generation
		// like JSON, XML, etc.
		c.MethodsCode.WriteString("\nfunc (m *" + currentModel.Name + ") SerializeBinary() ([]byte, error) {\n")

		// handle m is nil by returning []byte(0) and nil
		c.MethodsCode.WriteString("\tif m == nil {\n")
		c.MethodsCode.WriteString("\t\treturn []byte{0}, nil\n")
		c.MethodsCode.WriteString("\t}\n\n")

		c.MethodsCode.WriteString("\tbuf := new(bytes.Buffer)\n\n")
		for _, field := range currentModel.Fields {
			println("serialize current field: ", field.Name)
			isCustomType := c.Options.CCLDefinition.GetModelByName(field.Type) != nil
			switch field.Type {
			case cclValues.TypeNameString:
				c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(m." + field.Name + "))); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn nil, err\n")
				c.MethodsCode.WriteString("\t}\n")
				c.MethodsCode.WriteString("\tif _, err := buf.WriteString(m." + field.Name + "); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn nil, err\n")
				c.MethodsCode.WriteString("\t}\n")
			case cclValues.TypeNameBytes:
				c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(m." + field.Name + "))); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn nil, err\n")
				c.MethodsCode.WriteString("\t}\n")
				c.MethodsCode.WriteString("\tif _, err := buf.Write(m." + field.Name + "); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn nil, err\n")
				c.MethodsCode.WriteString("\t}\n")
			case cclValues.TypeNameDateTime:
				c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, m." + field.Name + ".UnixNano()); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn nil, err\n")
				c.MethodsCode.WriteString("\t}\n")
			default:
				if field.IsArray() {
					c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(m." + field.Name + "))); err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn nil, err\n")
					c.MethodsCode.WriteString("\t}\n")
					c.MethodsCode.WriteString("\tfor _, elem := range m." + field.Name + " {\n")
					if isCustomType {
						currentBytesName := "current_" + field.Name + "Bytes"
						c.MethodsCode.WriteString("\t\t" + currentBytesName + ", err := elem.SerializeBinary()\n")
						c.MethodsCode.WriteString("\t\tif err != nil {\n")
						c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
						c.MethodsCode.WriteString("\t\t}\n")
						c.MethodsCode.WriteString("\t\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(" + currentBytesName + "))); err != nil {\n")
						c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
						c.MethodsCode.WriteString("\t\t}\n")
						c.MethodsCode.WriteString("\t\tif err := binary.Write(buf, binary.LittleEndian, " + currentBytesName + "); err != nil {\n")
						c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
						c.MethodsCode.WriteString("\t\t}\n")
					} else {
						c.MethodsCode.WriteString("\t\tif err := binary.Write(buf, binary.LittleEndian, elem); err != nil {\n")
						c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
						c.MethodsCode.WriteString("\t\t}\n")
					}
					c.MethodsCode.WriteString("\t}\n")
				} else if isCustomType {
					currentBytesName := "current_" + field.Name + "Bytes"
					c.MethodsCode.WriteString("\t" + currentBytesName + ", err := m." + field.Name + ".SerializeBinary()\n")
					c.MethodsCode.WriteString("\tif err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn nil, err\n")
					c.MethodsCode.WriteString("\t}\n")
					c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(" + currentBytesName + "))); err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn nil, err\n")
					c.MethodsCode.WriteString("\t}\n")
					c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, " + currentBytesName + "); err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn nil, err\n")
					c.MethodsCode.WriteString("\t}\n")
				} else {
					c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, m." + field.Name + "); err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn nil, err\n")
					c.MethodsCode.WriteString("\t}\n")
				}
			}
		}

		c.MethodsCode.WriteString("\treturn buf.Bytes(), nil\n")
		c.MethodsCode.WriteString("}\n\n")

		c.MethodsCode.WriteString("func (m *" + currentModel.Name + ") DeserializeBinary(data []byte) error {\n")
		c.MethodsCode.WriteString("\tbuf := bytes.NewReader(data)\n\n")

		for _, field := range currentModel.Fields {
			isCustomType := c.Options.CCLDefinition.IsCustomType(field.Type)
			isPointer := isCustomType //TODO: Find a way to specify this
			fName := strings.ToLower(string(field.Name[0])) + field.Name[1:]
			fLenName := fName + "Len"
			fNameStrBytes := fName + "StrBytes"
			fNameUnix := fName + "Unix"
			fieldRealType := field.Type
			if isPointer {
				fieldRealType = "*" + fieldRealType
			}
			println("deserialize current field: ", field.Name)

			switch field.Type {
			case cclValues.TypeNameString:
				c.MethodsCode.WriteString("\tvar " + fLenName + " uint32\n")
				c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + fLenName + "); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn err\n")
				c.MethodsCode.WriteString("\t}\n")
				c.MethodsCode.WriteString("\t" + fNameStrBytes + " := make([]byte, " + fLenName + ")\n")
				c.MethodsCode.WriteString("\tif _, err := buf.Read(" + fNameStrBytes + "); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn err\n")
				c.MethodsCode.WriteString("\t}\n")
				c.MethodsCode.WriteString("\tm." + field.Name + " = string(" + fNameStrBytes + ")\n")
			case cclValues.TypeNameBytes:
				c.MethodsCode.WriteString("\tvar " + fLenName + " uint32\n")
				c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + fLenName + "); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn err\n")
				c.MethodsCode.WriteString("\t}\n")
				c.MethodsCode.WriteString("\tbytesData := make([]byte, " + fLenName + ")\n")
				c.MethodsCode.WriteString("\tif _, err := buf.Read(bytesData); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn err\n")
				c.MethodsCode.WriteString("\t}\n")
				c.MethodsCode.WriteString("\tm." + field.Name + " = bytesData\n")
			case cclValues.TypeNameDateTime:
				c.MethodsCode.WriteString("\tvar " + fNameUnix + " int64\n")
				c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + fNameUnix + "); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn err\n")
				c.MethodsCode.WriteString("\t}\n")
				c.MethodsCode.WriteString("\tm." + field.Name + " = time.Unix(0, " + fNameUnix + ")\n")
			default:
				if field.IsArray() {
					c.MethodsCode.WriteString("\tvar " + fLenName + " uint32\n")
					c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + fLenName + "); err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn err\n")
					c.MethodsCode.WriteString("\t}\n")
					c.MethodsCode.WriteString("\tm." + field.Name + " = make([]" + fieldRealType + ", " + fLenName + ")\n")
					c.MethodsCode.WriteString("\tfor i := uint32(0); i < " + fLenName + "; i++ {\n")
					if isCustomType {
						c.MethodsCode.WriteString("\t\tvar elem " + fieldRealType)
						if isPointer {
							c.MethodsCode.WriteString(" = new(" + field.Type + ")\n")
						} else {
							c.MethodsCode.WriteString("\n")
						}
						// we need to read the bytes from the buffer
						// and then deserialize the element
						c.MethodsCode.WriteString("\t\tvar elemLen uint32\n")
						c.MethodsCode.WriteString("\t\tif err := binary.Read(buf, binary.LittleEndian, &elemLen); err != nil {\n")
						c.MethodsCode.WriteString("\t\t\treturn err\n")
						c.MethodsCode.WriteString("\t\t}\n")
						c.MethodsCode.WriteString("\t\telemBytes := make([]byte, elemLen)\n")
						c.MethodsCode.WriteString("\t\tif _, err := buf.Read(elemBytes); err != nil {\n")
						c.MethodsCode.WriteString("\t\t\treturn err\n")
						c.MethodsCode.WriteString("\t\t}\n")
						c.MethodsCode.WriteString("\t\tif err := elem.DeserializeBinary(elemBytes); err != nil {\n")
						c.MethodsCode.WriteString("\t\t\treturn err\n")
						c.MethodsCode.WriteString("\t\t}\n")
						c.MethodsCode.WriteString("\t\tm." + field.Name + "[i] = elem\n")
					} else {
						c.MethodsCode.WriteString("\t\tif err := binary.Read(buf, binary.LittleEndian, &m." + field.Name + "[i]); err != nil {\n")
						c.MethodsCode.WriteString("\t\t\treturn err\n")
						c.MethodsCode.WriteString("\t\t}\n")
					}
					c.MethodsCode.WriteString("\t}\n")
				} else if isCustomType {
					// read the length of the next buffer that we need
					// var basicLen uint32
					// if err := binary.Read(buf, binary.LittleEndian, &basicLen); err != nil {
					// 	return err
					// }
					lenVarName := field.Name + "_bytesLen"
					c.MethodsCode.WriteString("\tvar " + lenVarName + " uint32\n")
					c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + lenVarName + "); err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn err\n")
					c.MethodsCode.WriteString("\t}\n")

					bytesVarName := field.Name + "Bytes"
					c.MethodsCode.WriteString("\t" + bytesVarName + " := make([]byte, " + lenVarName + ")\n")
					c.MethodsCode.WriteString("\tif _, err := buf.Read(" + bytesVarName + "); err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn err\n")
					c.MethodsCode.WriteString("\t}\n")

					// make sure m.field is not nil
					c.MethodsCode.WriteString("\tif m." + field.Name + " == nil {\n")
					c.MethodsCode.WriteString("\t\tm." + field.Name + " = new(" + field.Type + ")\n")
					c.MethodsCode.WriteString("\t}\n")

					c.MethodsCode.WriteString("\tif err := m." + field.Name + ".DeserializeBinary(" + bytesVarName + "); err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn err\n")
					c.MethodsCode.WriteString("\t}\n")
				} else {
					c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &m." + field.Name + "); err != nil {\n")
					c.MethodsCode.WriteString("\t\treturn err\n")
					c.MethodsCode.WriteString("\t}\n")
				}
			}
		}

		c.MethodsCode.WriteString("\treturn nil\n")
		c.MethodsCode.WriteString("}\n\n")

	}

	return nil
}
