package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"

	"github.com/diamondburned/cchat/repository"
)

var cchatPkg = repository.Main[repository.RootPath]

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage:", filepath.Base(os.Args[0]), "<out_path.go>")
	}

	var (
		buf bytes.Buffer
		ent bytes.Buffer
		con bytes.Buffer
		val bytes.Buffer
	)

	buf.WriteString("// Code generated by entity_types. DO NOT EDIT.\n\n")
	buf.WriteString("package netrpc\n\n")
	buf.WriteString(`import "github.com/diamondburned/cchat"` + "\n\n")
	buf.WriteString("const (\n")

	for _, iface := range cchatPkg.Interfaces {
		if !isIdentifier(&iface) {
			continue
		}

		fmt.Fprintf(&ent, "\t%[1]sEntity EntityType = %[1]q\n", iface.Name)
		fmt.Fprintf(&con, "\tcase cchat.%[1]s: return %[1]sEntity\n", iface.Name)
		fmt.Fprintf(&val, "\n\t\t%sEntity,", iface.Name)
	}

	buf.Write(ent.Bytes())
	buf.WriteString(")\n\n")

	buf.WriteString("// QueryEntityType resolves the entity type for the given value.\n")
	buf.WriteString("// An empty string is returned if the value is unknown.\n")
	buf.WriteString("func QueryEntityType(v cchat.Identifier) EntityType {\n")

	buf.WriteString("\tswitch v.(type) {\n")
	buf.Write(con.Bytes())
	buf.WriteString("\t}\n\n")
	buf.WriteString("\treturn \"\"\n")
	buf.WriteString("}\n\n")

	buf.WriteString("// IsValid validates that the entity type is valid.\n")
	buf.WriteString("func (t EntityType) IsValid() bool {")
	buf.WriteString("\tswitch t {\n")
	buf.WriteString("\tcase")
	buf.Write(bytes.TrimSuffix(val.Bytes(), []byte(",")))
	buf.WriteString(":\n")
	buf.WriteString("\t\treturn true\n")
	buf.WriteString("\tdefault:")
	buf.WriteString("\t\treturn false\n")
	buf.WriteString("\t}")
	buf.WriteString("}")

	fmt, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalln("invalid code:", err)
	}

	if err := os.WriteFile(os.Args[1], fmt, 0666); err != nil {
		log.Fatalln("failed to write file:", err)
	}
}

func isIdentifier(iface *repository.Interface) bool {
	for _, embed := range iface.Embeds {
		if embed.InterfaceName == "Identifier" {
			return true
		}

		if iface := cchatPkg.Interface(embed.InterfaceName); iface != nil {
			if isIdentifier(iface) {
				return true
			}
		}
	}
	return false
}
