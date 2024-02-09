package common

// this module exists to avoid circular references between the log and config packages
type EnvGetter func(string) (string, bool)
