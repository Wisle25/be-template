package generator

// IdGenerator interface defines a method for generating unique identifiers.
// This interface is typically implemented to generate IDs using UUID version 7 or any other ID generation strategy.
type IdGenerator interface {
	// Generate generates a new unique identifier and returns it as a string.
	Generate() string
}
