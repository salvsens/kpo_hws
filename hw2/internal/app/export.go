package app

// ResultExporter writes a calculated result to a file.
type ResultExporter interface {
	// WriteResult saves the result at the supplied path.
	WriteResult(path string, result Result) error
}

// ExportResult calculates the result and saves it at the supplied path.
func (g *Gradebook) ExportResult(studentID string, mode ResultMode, path string) error {
	result, err := g.CalculateResult(studentID, mode)
	if err != nil {
		return err
	}
	return g.exporter.WriteResult(path, result)
}
