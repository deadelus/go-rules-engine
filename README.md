# 🚀 Go Rules Engine

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Test Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://github.com/deadelus/go-rules-engine)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Un moteur de règles métier puissant et flexible pour Go, inspiré de [json-rules-engine](https://github.com/CacheControl/json-rules-engine). Évaluez des conditions complexes et déclenchez des événements basés sur des faits dynamiques.

## ✨ Fonctionnalités

- 🎯 **Règles définies en JSON ou en code** - Chargez vos règles depuis des fichiers JSON ou créez-les directement en Go
- 🔄 **Conditions complexes** - Supportez les opérateurs `all` et `any` avec imbrication infinie
- 📊 **Opérateurs riches** - `equal`, `not_equal`, `greater_than`, `less_than`, `in`, `not_in`, `contains`, `not_contains`
- 🎪 **Système d'événements** - Callbacks personnalisés et handlers globaux pour réagir aux résultats
- 💾 **Faits dynamiques** - Calculez des valeurs à la volée avec des callbacks
- 🧮 **Support JSONPath** - Accédez à des données imbriquées avec `$.path.to.value`
- ⚡ **Priorités de règles** - Contrôlez l'ordre d'évaluation avec des priorités
- 🔒 **Thread-safe** - Protégé par des mutex pour un usage concurrent
- ✅ **100% de couverture de tests** - Code robuste et testé en profondeur

## 📦 Installation

```bash
go get github.com/deadelus/go-rules-engine
```

## 🚀 Démarrage rapide

### Exemple basique

```go
package main

import (
    "fmt"
    gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
    // 1. Créer le moteur de règles
    engine := gorulesengine.NewEngine()

    // 2. Définir une règle
    rule := &gorulesengine.Rule{
        Name:     "adult-user",
        Priority: 10,
        Conditions: gorulesengine.ConditionSet{
            All: []gorulesengine.ConditionNode{
                {
                    Condition: &gorulesengine.Condition{
                        Fact:     "age",
                        Operator: "greater_than",
                        Value:    18,
                    },
                },
            },
        },
        Event: gorulesengine.Event{
            Type: "user-is-adult",
            Params: map[string]interface{}{
                "message": "Utilisateur majeur détecté",
            },
        },
    }

    // 3. Ajouter la règle au moteur
    engine.AddRule(rule)

    // 4. Créer l'almanac avec des faits
    almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
    almanac.AddFact("age", 25)

    // 5. Exécuter le moteur
    results, err := engine.Run(almanac)
    if err != nil {
        panic(err)
    }

    // 6. Afficher les résultats
    for _, result := range results {
        if result.Result {
            fmt.Printf("✅ Règle '%s' déclenchée!\n", result.Rule.Name)
            fmt.Printf("   Event: %s\n", result.Event.Type)
        }
    }
}
```

### Charger des règles depuis JSON

```go
package main

import (
    "encoding/json"
    "fmt"
    gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
    // JSON de la règle
    ruleJSON := `{
        "name": "premium-user",
        "priority": 10,
        "conditions": {
            "all": [
                {
                    "condition": {
                        "fact": "accountType",
                        "operator": "equal",
                        "value": "premium"
                    }
                },
                {
                    "condition": {
                        "fact": "revenue",
                        "operator": "greater_than",
                        "value": 1000
                    }
                }
            ]
        },
        "event": {
            "type": "premium-user-detected",
            "params": {
                "discount": 20
            }
        }
    }`

    var rule gorulesengine.Rule
    json.Unmarshal([]byte(ruleJSON), &rule)

    engine := gorulesengine.NewEngine()
    engine.AddRule(&rule)

    almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
    almanac.AddFact("accountType", "premium")
    almanac.AddFact("revenue", 1500)

    results, _ := engine.Run(almanac)
    fmt.Printf("Règles déclenchées: %d\n", len(results))
}
```

### Charger des règles ET des facts depuis JSON

```go
package main

import (
    "encoding/json"
    "fmt"
    gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
    // JSON des règles
    rulesJSON := `[
        {
            "name": "high-value-order",
            "priority": 100,
            "conditions": {
                "all": [
                    {
                        "condition": {
                            "fact": "user.isPremium",
                            "operator": "equal",
                            "value": true
                        }
                    },
                    {
                        "condition": {
                            "fact": "order.total",
                            "operator": "greater_than",
                            "value": 100
                        }
                    }
                ]
            },
            "event": {
                "type": "premium-discount",
                "params": {"discount": 25}
            }
        }
    ]`

    // JSON des facts (données)
    factsJSON := `{
        "user": {
            "id": 12345,
            "isPremium": true,
            "name": "Alice"
        },
        "order": {
            "id": "ORD-001",
            "total": 150.50
        }
    }`

    // Charger les règles
    var rules []*gorulesengine.Rule
    json.Unmarshal([]byte(rulesJSON), &rules)

    // Charger les facts
    var factsData map[string]interface{}
    json.Unmarshal([]byte(factsJSON), &factsData)

    // Créer l'engine et ajouter les règles
    engine := gorulesengine.NewEngine()
    for _, rule := range rules {
        engine.AddRule(rule)
    }

    // Créer l'almanac et ajouter les facts
    almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
    for key, value := range factsData {
        almanac.AddFact(gorulesengine.FactID(key), value)
    }

    // Exécuter
    results, _ := engine.Run(almanac)
    fmt.Printf("Règles déclenchées: %d\n", len(results))
}
```

## 📖 Documentation

### Architecture

Le moteur de règles est composé de plusieurs composants clés :

#### 1. **Engine** - Le moteur principal

```go
engine := gorulesengine.NewEngine()
engine.AddRule(rule)
results, err := engine.Run(almanac)
```

**Méthodes :**
- `AddRule(rule *Rule)` - Ajoute une règle au moteur
- `Run(almanac *Almanac) ([]RuleResult, error)` - Exécute toutes les règles
- `RegisterCallback(name string, callback Callback)` - Enregistre un callback nommé
- `On(outcome string, handler EventHandler)` - Handler global pour success/failure
- `OnEvent(eventType string, handler EventHandler)` - Handler spécifique à un type d'événement

#### 2. **Rule** - Une règle métier

```go
rule := &gorulesengine.Rule{
    Name:       "my-rule",
    Priority:   10,          // Plus élevé = exécuté en premier
    Conditions: conditionSet,
    Event:      event,
    OnSuccess:  strPtr("mySuccessCallback"), // Optionnel
    OnFailure:  strPtr("myFailureCallback"), // Optionnel
}
```

#### 3. **Condition** - Une condition à évaluer

```go
condition := &gorulesengine.Condition{
    Fact:     "age",
    Operator: "greater_than",
    Value:    18,
    Path:     "$.user.age", // Optionnel: JSONPath pour données imbriquées
}
```

**Opérateurs disponibles :**
- `equal` - Égalité
- `not_equal` - Différent de
- `greater_than` - Supérieur à
- `greater_than_or_equal` - Supérieur ou égal à
- `less_than` - Inférieur à
- `less_than_or_equal` - Inférieur ou égal à
- `in` - Dans la liste
- `not_in` - Pas dans la liste
- `contains` - Contient (pour strings et arrays)
- `not_contains` - Ne contient pas

#### 4. **ConditionSet** - Groupement de conditions

```go
// Toutes les conditions doivent être vraies (AND)
conditionSet := gorulesengine.ConditionSet{
    All: []gorulesengine.ConditionNode{
        {Condition: &condition1},
        {Condition: &condition2},
    },
}

// Au moins une condition doit être vraie (OR)
conditionSet := gorulesengine.ConditionSet{
    Any: []gorulesengine.ConditionNode{
        {Condition: &condition1},
        {Condition: &condition2},
    },
}

// Imbrication (AND de OR)
conditionSet := gorulesengine.ConditionSet{
    All: []gorulesengine.ConditionNode{
        {Condition: &condition1},
        {
            ConditionSet: &gorulesengine.ConditionSet{
                Any: []gorulesengine.ConditionNode{
                    {Condition: &condition2},
                    {Condition: &condition3},
                },
            },
        },
    },
}
```

#### 5. **Almanac** - Stockage des faits

```go
almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

// Ajouter des faits simples
almanac.AddFact("age", 25)
almanac.AddFact("country", "FR")

// Ajouter des faits dynamiques
almanac.AddFact("temperature", gorulesengine.Fact{
    ID: "temperature",
    Calculate: func(params map[string]interface{}, almanac *gorulesengine.Almanac) (interface{}, error) {
        // Logique de calcul personnalisée
        return fetchTemperature(), nil
    },
})

// Récupérer un fait
value, err := almanac.GetFactValue("age", nil)
```

#### 6. **Event** - Événement déclenché

```go
event := gorulesengine.Event{
    Type: "user-approved",
    Params: map[string]interface{}{
        "userId": 123,
        "reason": "All conditions met",
    },
}
```

### Système de callbacks et handlers

#### Callbacks nommés (définis dans les règles JSON)

```go
engine := gorulesengine.NewEngine()

// Enregistrer le callback
engine.RegisterCallback("sendEmail", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
    fmt.Printf("Envoi d'email pour: %s\n", event.Type)
    return nil
})

// Dans la règle JSON
rule := &gorulesengine.Rule{
    Name: "email-rule",
    OnSuccess: strPtr("sendEmail"), // Référence au callback
    // ...
}
```

#### Handlers globaux

```go
// Handler pour toutes les règles réussies
engine.On("success", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
    fmt.Printf("✅ Règle réussie: %s\n", ruleResult.Rule.Name)
    return nil
})

// Handler pour toutes les règles échouées
engine.On("failure", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
    fmt.Printf("❌ Règle échouée: %s\n", ruleResult.Rule.Name)
    return nil
})
```

#### Handlers par type d'événement

```go
// Handler spécifique pour un type d'événement
engine.OnEvent("user-approved", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
    userId := event.Params["userId"]
    fmt.Printf("Utilisateur %v approuvé!\n", userId)
    return nil
})
```

### Support JSONPath

Accédez à des données imbriquées dans vos faits :

```go
almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
almanac.AddFact("user", map[string]interface{}{
    "profile": map[string]interface{}{
        "age": 25,
        "address": map[string]interface{}{
            "city": "Paris",
        },
    },
})

// Utilisez JSONPath dans les conditions
condition := &gorulesengine.Condition{
    Fact:     "user",
    Path:     "$.profile.address.city",
    Operator: "equal",
    Value:    "Paris",
}
```

### Gestion des erreurs

Le moteur utilise un système d'erreurs typées pour une meilleure traçabilité :

```go
results, err := engine.Run(almanac)
if err != nil {
    var ruleErr *gorulesengine.RuleEngineError
    if errors.As(err, &ruleErr) {
        fmt.Printf("Type: %s, Message: %s\n", ruleErr.Type, ruleErr.Msg)
    }
}
```

**Types d'erreurs :**
- `ErrEngine` - Erreur générale du moteur
- `ErrAlmanac` - Erreur liée aux faits (almanac)
- `ErrFact` - Erreur de calcul de fait
- `ErrRule` - Erreur dans la définition de la règle
- `ErrCondition` - Erreur d'évaluation de condition
- `ErrOperator` - Opérateur invalide ou non trouvé
- `ErrEvent` - Erreur liée aux événements
- `ErrJSON` - Erreur de parsing JSON

## 🧪 Tests

Le projet dispose d'une couverture de tests de **100%** :

```bash
# Exécuter tous les tests
go test ./src -v

# Avec couverture
go test ./src -coverprofile=coverage.out
go tool cover -html=coverage.out

# Voir le résumé
go tool cover -func=coverage.out | tail -1
# Output: total: (statements) 100.0%
```

## 🔍 Qualité du code

Le code respecte toutes les conventions Go et passe les linters sans avertissement :

```bash
# go vet (vérification statique)
go vet ./src/...

# golint (style Go)
golint ./src/...

# Format du code
go fmt ./src/...
```

**Standards respectés:**
- ✅ Conventions de nommage Go (CamelCase, pas de ALL_CAPS)
- ✅ Documentation GoDoc complète sur toutes les exports
- ✅ Gestion d'erreurs appropriée
- ✅ Code thread-safe avec mutexes
- ✅ Tests exhaustifs avec 100% de couverture

## 🗺️ Roadmap

### ✅ Phases complétées

- [x] Phase 1: Structures de base (Condition, Rule, Fact)
- [x] Phase 2: Almanac et gestion des faits
- [x] Phase 3: Opérateurs (equal, greater_than, less_than, etc.)
- [x] Phase 4: Évaluation des conditions (all/any, imbrication)
- [x] Phase 5: Engine avec système d'événements
- [x] Phase 6: Support JSON et désérialisation
- [x] Phase 7: Features avancées (callbacks, handlers, JSONPath)
- [x] Tests complets avec 100% de couverture

### 🚧 Phases à venir

#### Phase 8: API ergonomique et builders

**Builders fluent pour créer des règles**
```go
rule := NewRuleBuilder().
    WithName("adult-user").
    WithPriority(10).
    WithCondition(Equal("age", 18)).
    WithEvent("user-is-adult", nil).
    Build()
```

**Helpers de conditions**
```go
condition := All(
    GreaterThan("age", 18),
    Equal("country", "FR"),
    Any(
        Equal("status", "premium"),
        Equal("status", "vip"),
    ),
)
```

#### Phase 9: Documentation et exemples

- [x] Documentation GoDoc complète
- [x] Exemples dans `examples/`
  - [x] `examples/full-demo.go` - Démonstration complète de toutes les fonctionnalités
  - [x] `examples/basic/` - Cas simple
  - [x] `examples/json/` - Chargement JSON
  - [x] `examples/advanced/` - Features avancées
  - [x] `examples/custom-operator/` - Opérateurs personnalisés

#### Phase 10: Nouveaux opérateurs

- [ ] `regex` - Vérifier si la valeur correspond à une expression régulière

#### Phase 11: Performance et optimisation

- [ ] Benchmarks complets
- [ ] Cache des résultats de conditions
- [ ] Évaluation parallèle des règles indépendantes
- [ ] Profilage mémoire et CPU

#### Phase 12: Features avancées

- [ ] Tri des fact par `priority`
- [ ] Support de règles async
- [ ] Persistance des résultats
- [ ] Métriques et monitoring
- [ ] Hot-reload des règles
- [ ] API REST optionnelle

## 🤝 Contribution

Les contributions sont les bienvenues ! Pour contribuer :

1. Forkez le projet
2. Créez une branche (`git checkout -b feature/amazing-feature`)
3. Committez vos changements (`git commit -m 'Add amazing feature'`)
4. Pushez vers la branche (`git push origin feature/amazing-feature`)
5. Ouvrez une Pull Request

**Guidelines :**
- Écrivez des tests pour toutes les nouvelles fonctionnalités
- Maintenez la couverture à 100%
- Suivez les conventions Go (gofmt, golint)
- Documentez vos fonctions publiques

## 📄 License

Ce projet est sous licence MIT. Voir le fichier [LICENSE](LICENSE) pour plus de détails.

**Copyright (c) 2026 Geoffrey Trambolho (@deadelus)**

## 🙏 Remerciements

Inspiré par [json-rules-engine](https://github.com/CacheControl/json-rules-engine) de CacheControl.

## 📞 Contact

Créé par [@deadelus](https://github.com/deadelus)

---

⭐ N'oubliez pas de donner une étoile si ce projet vous aide !
