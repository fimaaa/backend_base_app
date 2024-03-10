package helper

func  (u HTTPHelper) RemoveIndex(s []string, index int) []string {
    return append(s[:index], s[index:]...)
}