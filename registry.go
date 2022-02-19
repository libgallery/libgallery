package libgallery

// The registry contains a list of names suitable for
// identification mapped to implementations.
// The key should not have spaces or captialisation.
var Registry map[string]Driver = map[string]Driver{}

// Register new implementation, should
// only be accessed during Init() by
// implementation packages!
func Register(name string, impl Driver) {
	Registry[name] = impl
}
