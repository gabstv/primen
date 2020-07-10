// The MIT License (MIT)

// Copyright (c) 2015 Aymerick JEHANNE

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package css

import "fmt"

// Declaration represents a parsed style property
type Declaration struct {
	Property  string
	Value     string
	Important bool
}

// NewDeclaration instanciates a new Declaration
func NewDeclaration() *Declaration {
	return &Declaration{}
}

// Returns string representation of the Declaration
func (decl *Declaration) String() string {
	return decl.StringWithImportant(true)
}

// StringWithImportant returns string representation with optional !important part
func (decl *Declaration) StringWithImportant(option bool) string {
	result := fmt.Sprintf("%s: %s", decl.Property, decl.Value)

	if option && decl.Important {
		result += " !important"
	}

	result += ";"

	return result
}

// Equal returns true if both Declarations are equals
func (decl *Declaration) Equal(other *Declaration) bool {
	return (decl.Property == other.Property) && (decl.Value == other.Value) && (decl.Important == other.Important)
}

//
// DeclarationsByProperty
//

// DeclarationsByProperty represents sortable style declarations
type DeclarationsByProperty []*Declaration

// Implements sort.Interface
func (declarations DeclarationsByProperty) Len() int {
	return len(declarations)
}

// Implements sort.Interface
func (declarations DeclarationsByProperty) Swap(i, j int) {
	declarations[i], declarations[j] = declarations[j], declarations[i]
}

// Implements sort.Interface
func (declarations DeclarationsByProperty) Less(i, j int) bool {
	return declarations[i].Property < declarations[j].Property
}
