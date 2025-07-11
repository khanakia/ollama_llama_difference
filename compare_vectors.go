package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type LlamaEntry struct {
	Index     int         `json:"index"`
	Embedding [][]float64 `json:"embedding"`
}

type OllamaResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float64 `json:"embeddings"`
}

func normalizeVector(vec []float64) []float64 {
	var sum float64
	for _, v := range vec {
		sum += v * v
	}

	norm := 0.0
	if sum > 0 {
		norm = 1.0 / math.Sqrt(sum)
	}

	result := make([]float64, len(vec))
	for i := range vec {
		result[i] = vec[i] * norm
	}
	return result
}

func calculateDifference(vec1, vec2 []float64) float64 {
	if len(vec1) != len(vec2) {
		return -1 // Error indicator
	}

	var sum float64
	for i := range vec1 {
		diff := vec1[i] - vec2[i]
		sum += diff * diff
	}

	return math.Sqrt(sum) / float64(len(vec1)) * 100 // Return as percentage
}

func compareVectors() {
	// Read llama.json
	llamaData, err := os.ReadFile("llama.json")
	if err != nil {
		fmt.Printf("Error reading llama.json: %v\n", err)
		return
	}

	var llamaEntries []LlamaEntry
	err = json.Unmarshal(llamaData, &llamaEntries)
	if err != nil {
		fmt.Printf("Error parsing llama.json: %v\n", err)
		return
	}

	// Read ollama.json
	ollamaData, err := os.ReadFile("ollama.json")
	if err != nil {
		fmt.Printf("Error reading ollama.json: %v\n", err)
		return
	}

	var ollamaResponse OllamaResponse
	err = json.Unmarshal(ollamaData, &ollamaResponse)
	if err != nil {
		fmt.Printf("Error parsing ollama.json: %v\n", err)
		return
	}

	// Generate markdown report
	report := "# Vector Comparison Report\n\n"
	report += "## Overview\n\n"
	report += fmt.Sprintf("- **Llama vectors count**: %d\n", len(llamaEntries))
	report += fmt.Sprintf("- **Ollama vectors count**: %d\n", len(ollamaResponse.Embeddings))
	report += fmt.Sprintf("- **Vector dimensions**: %d\n", len(ollamaResponse.Embeddings[0]))
	report += "\n"

	// Compare vectors
	report += "## Vector Comparison Results\n\n"
	report += "| Vector Index | Llama (Normalized) | Ollama | Difference (%) |\n"
	report += "|--------------|-------------------|--------|----------------|\n"

	totalDiff := 0.0
	validComparisons := 0

	for i := 0; i < len(llamaEntries) && i < len(ollamaResponse.Embeddings); i++ {
		// Normalize llama vector
		normalizedLlama := normalizeVector(llamaEntries[i].Embedding[0])

		fmt.Println(normalizedLlama)

		// Calculate difference
		diff := calculateDifference(normalizedLlama, ollamaResponse.Embeddings[i])

		if diff >= 0 {
			totalDiff += diff
			validComparisons++

			// Show first few values as sample
			llamaSample := fmt.Sprintf("[%.6f, %.6f, %.6f, ...]",
				normalizedLlama[0], normalizedLlama[1], normalizedLlama[2])
			ollamaSample := fmt.Sprintf("[%.6f, %.6f, %.6f, ...]",
				ollamaResponse.Embeddings[i][0], ollamaResponse.Embeddings[i][1], ollamaResponse.Embeddings[i][2])

			report += fmt.Sprintf("| %d | %s | %s | %.4f%% |\n",
				i, llamaSample, ollamaSample, diff)
		}
	}

	avgDiff := 0.0
	if validComparisons > 0 {
		avgDiff = totalDiff / float64(validComparisons)
	}

	report += "\n## Summary Statistics\n\n"
	report += fmt.Sprintf("- **Total vectors compared**: %d\n", validComparisons)
	report += fmt.Sprintf("- **Average difference**: %.4f%%\n", avgDiff)
	report += fmt.Sprintf("- **Total difference**: %.4f%%\n", totalDiff)

	// Write report to file
	err = os.WriteFile("vector_comparison_report.md", []byte(report), 0644)
	if err != nil {
		fmt.Printf("Error writing report: %v\n", err)
		return
	}

	fmt.Println("Vector comparison report generated: vector_comparison_report.md")
	fmt.Printf("Average difference: %.4f%%\n", avgDiff)
}
