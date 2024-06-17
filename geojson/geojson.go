package geojson

import (
	"encoding/json"
	"fmt"

	"github.com/twpayne/go-geos/geometry"
)

const (
	featureType           = "Feature"
	featureCollectionType = "FeatureCollection"
)

// A Feature is a feature.
type Feature struct {
	ID         interface{}
	Geometry   geometry.Geometry
	Properties map[string]interface{}
}

// A FeatureCollection is a feature collection.
type FeatureCollection struct {
	Features   []*Feature
	Properties map[string]interface{}
}

type feature struct {
	ID         interface{}            `json:"id,omitempty"`
	Type       string                 `json:"type"`
	Geometry   *geometry.Geometry     `json:"geometry"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type featureCollection struct {
	Type       string                 `json:"type"`
	Features   []feature              `json:"features"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (f *Feature) MarshalJSON() ([]byte, error) {
	return json.Marshal(feature{
		ID:         f.ID,
		Type:       featureType,
		Geometry:   &f.Geometry,
		Properties: f.Properties,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Feature) UnmarshalJSON(data []byte) error {
	var geoJSONFeature feature
	if err := json.Unmarshal(data, &geoJSONFeature); err != nil {
		return err
	}
	if geoJSONFeature.Type != featureType {
		return fmt.Errorf("not a Feature: %s", geoJSONFeature.Type)
	}
	f.ID = geoJSONFeature.ID
	f.Geometry = *geoJSONFeature.Geometry
	f.Properties = geoJSONFeature.Properties
	return nil
}

// MarshalJSON implements json.Marshaler.
func (fc FeatureCollection) MarshalJSON() ([]byte, error) {
	features := make([]feature, 0, len(fc.Features))
	for _, f := range fc.Features {
		feature := feature{
			ID:         f.ID,
			Type:       featureType,
			Geometry:   &f.Geometry,
			Properties: f.Properties,
		}
		features = append(features, feature)
	}
	return json.Marshal(featureCollection{
		Type:       featureCollectionType,
		Properties: fc.Properties,
		Features:   features,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (fc *FeatureCollection) UnmarshalJSON(data []byte) error {
	var geoJSONFeatureCollection featureCollection
	if err := json.Unmarshal(data, &geoJSONFeatureCollection); err != nil {
		return err
	}
	if geoJSONFeatureCollection.Type != featureCollectionType {
		return fmt.Errorf("not a FeatureCollection: %s", geoJSONFeatureCollection.Type)
	}
	featureCollection := make([]*Feature, 0, len(geoJSONFeatureCollection.Features))
	for _, feature := range geoJSONFeatureCollection.Features {
		if feature.Type != featureType {
			return fmt.Errorf("not a Feature: %s", feature.Type)
		}
		f := &Feature{
			ID:         feature.ID,
			Geometry:   *feature.Geometry,
			Properties: feature.Properties,
		}
		featureCollection = append(featureCollection, f)
	}
	fc.Features = featureCollection
	fc.Properties = geoJSONFeatureCollection.Properties
	return nil
}
