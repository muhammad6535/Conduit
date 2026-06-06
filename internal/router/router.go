package router

import (
	"fmt"
	"regexp"

	"github.com/conduit-ai/gateway/internal/config"
	"github.com/conduit-ai/gateway/internal/provider"
)

type Route struct {
	Provider provider.Provider
	Model    string
}

type Router struct {
	providers map[string]provider.Provider
	cfg       *config.RoutingConfig
	rules     []compiledRule
}

type compiledRule struct {
	pattern  *regexp.Regexp
	provider string
	model    string
}

func New(providers map[string]provider.Provider, cfg *config.RoutingConfig) (*Router, error) {
	r := &Router{
		providers: providers,
		cfg:       cfg,
	}

	for _, rule := range cfg.Rules {
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, fmt.Errorf("compile rule pattern %q: %w", rule.Pattern, err)
		}
		r.rules = append(r.rules, compiledRule{
			pattern:  re,
			provider: rule.Provider,
			model:    rule.Model,
		})
	}

	return r, nil
}

func (r *Router) Route(modelName string) (*Route, error) {
	if modelName == "" {
		modelName = r.cfg.DefaultModel
	}

	for _, rule := range r.rules {
		if rule.pattern.MatchString(modelName) {
			prov, ok := r.providers[rule.provider]
			if !ok {
				return r.defaultRoute()
			}
			model := rule.model
			if model == "" {
				model = modelName
			}
			return &Route{Provider: prov, Model: model}, nil
		}
	}

	for _, prov := range r.providers {
		if hasModel(prov, modelName) {
			return &Route{Provider: prov, Model: modelName}, nil
		}
	}

	return r.defaultRoute()
}

func (r *Router) defaultRoute() (*Route, error) {
	prov, ok := r.providers[r.cfg.DefaultProvider]
	if !ok {
		for _, p := range r.providers {
			prov = p
			break
		}
	}
	if prov == nil {
		return nil, fmt.Errorf("no providers configured")
	}
	return &Route{Provider: prov, Model: r.cfg.DefaultModel}, nil
}

func (r *Router) AvailableModels() []map[string]string {
	var models []map[string]string
	for name, prov := range r.providers {
		models = append(models, map[string]string{
			"id":       name + "/*",
			"provider": prov.Name(),
			"object":   "model",
		})
	}
	return models
}

func hasModel(prov provider.Provider, model string) bool {
	// For now, assume the provider can handle the model if its name
	// contains a known prefix. This will be improved with a model registry.
	pName := prov.Name()
	switch pName {
	case "openai":
		return true // OpenAI can handle any OpenAI model
	case "anthropic":
		return true // Anthropic can handle any Anthropic model
	default:
		return true
	}
}
