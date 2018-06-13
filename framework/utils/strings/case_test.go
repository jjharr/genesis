/*
 * Inflector Pkg (Go)
 *
 * Copyright (c) 2013 Ivan Torres
 * Released under the MIT license
 * https://github.com/mexpolk/inflector/blob/master/LICENSE
 *
 */

package strings

import (
	"testing"
)

type inflectionSample struct {
	str, out string
}

func TestToCamel(t *testing.T) {
	samples := []inflectionSample{
		{"sample text", "sampleText"},
		{"sample-text", "sampleText"},
		{"sample_text", "sampleText"},
		{"sampleText", "sampleText"},
		{"sample 2 Text", "sample2Text"},
	}

	for _, sample := range samples {
		if out := ToCamel(sample.str); out != sample.out {
			t.Errorf("got %q, expected %q", out, sample.out)
		}
	}
}

func TestToDash(t *testing.T) {
	samples := []inflectionSample{
		{"sample text", "sample-text"},
		{"sample-text", "sample-text"},
		{"sample_text", "sample-text"},
		{"sampleText", "sample-text"},
		{"sample 2 Text", "sample-2-text"},
	}

	for _, sample := range samples {
		if out := ToDash(sample.str); out != sample.out {
			t.Errorf("got %q, expected %q", out, sample.out)
		}
	}
}

func TestToPascal(t *testing.T) {
	samples := []inflectionSample{
		{"sample text", "SampleText"},
		{"sample-text", "SampleText"},
		{"sample_text", "SampleText"},
		{"sampleText", "SampleText"},
		{"sample 2 Text", "Sample2Text"},
	}

	for _, sample := range samples {
		if out := ToPascal(sample.str); out != sample.out {
			t.Errorf("got %q, expected %q", out, sample.out)
		}
	}
}

func TestToSnake(t *testing.T) {
	samples := []inflectionSample{
		{"sample text", "sample_text"},
		{"sample-text", "sample_text"},
		{"sample_text", "sample_text"},
		{"sampleText", "sample_text"},
		{"sample 2 Text", "sample_2_text"},
	}

	for _, sample := range samples {
		if out := ToSnake(sample.str); out != sample.out {
			t.Errorf("got %q, expected %q", out, sample.out)
		}
	}
}

func TestToTitle(t *testing.T) {
	samples := []inflectionSample{
		{"sample text", "Sample Text"},
		{"sample-text", "Sample Text"},
		{"sample_text", "Sample Text"},
		{"sampleText", "Sample Text"},
		{"sample 2 Text", "Sample 2 Text"},
	}

	for _, sample := range samples {
		if out := ToTitle(sample.str); out != sample.out {
			t.Errorf("got %q, expected %q", out, sample.out)
		}
	}
}

func TestToHeader(t *testing.T) {
	samples := []inflectionSample{
		{"sample text", "Sample-Text"},
		{"sample-text", "Sample-Text"},
		{"sample_text", "Sample-Text"},
		{"sampleText", "Sample-Text"},
		{"sample 2 Text", "Sample-2-Text"},
	}

	for _, sample := range samples {
		if out := ToHeader(sample.str); out != sample.out {
			t.Errorf("got %q, expected %q", out, sample.out)
		}
	}
}
