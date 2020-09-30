package ibos

func missingFiles(lastFile string, filenames []string) []string {
	for i, filename := range filenames {
		if filename == lastFile {
			return filenames[i+1:]
		}
	}
	return filenames
}
