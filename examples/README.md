# Examples - Go Rules Engine

Ce dossier contient des exemples d'utilisation du moteur de règles Go.

## 📚 Exemples disponibles

### 1. basic/main.go
**Exemple basique** - Vérification d'âge simple avec une seule règle.

```bash
go run examples/basic/main.go
```

Démontre :
- ✅ Création d'un engine
- ✅ Règle simple avec condition
- ✅ Opérateur `greater_than`
- ✅ Tests avec différentes valeurs

### 2. json/main.go
**Chargement JSON** - Charge des règles et facts depuis JSON.

```bash
go run examples/json/main.go
```

Démontre :
- ✅ Unmarshal de règles JSON
- ✅ Unmarshal de facts JSON
- ✅ Ajout de règles à l'engine
- ✅ Ajout de facts à l'almanac
- ✅ Règles VIP et régulières

### 3. custom-operator/main.go
**Opérateurs personnalisés** - Création d'opérateurs custom.

```bash
go run examples/custom-operator/main.go
```

Démontre :
- ✅ Interface `Operator`
- ✅ Implémentation de `CustomOperator`
- ✅ Opérateurs `starts_with`, `ends_with`, `between`
- ✅ `RegisterOperator` pour enregistrer les opérateurs

### 4. advanced/main.go
**Fonctionnalités avancées** - Callbacks, handlers et dynamic facts.

```bash
go run examples/advanced/main.go
```

Démontre :
- ✅ Callbacks nommés avec `RegisterCallback`
- ✅ Handler global `OnSuccess`
- ✅ Handler spécifique par type d'événement `On()`
- ✅ Dynamic facts (calcul de remise)
- ✅ Multiple handlers simultanés

### 5. full-demo.go
**Démonstration complète** - Toutes les fonctionnalités en un seul exemple.

```bash
go run examples/full-demo.go
```

Démontre :
- ✅ Règles simples et complexes
- ✅ Conditions imbriquées (all/any)
- ✅ Callbacks et handlers
- ✅ Chargement JSON
- ✅ Dynamic facts
- ✅ JSONPath
- ✅ Historique des événements

## 🚀 Exécution

Depuis la racine du projet :

```bash
# Exemple basique
go run examples/basic/main.go

# JSON
go run examples/json/main.go

# Custom operators
go run examples/custom-operator/main.go

# Advanced
go run examples/advanced/main.go

# Full demo
go run examples/full-demo.go
```

## 📖 Documentation complète

Voir le [README principal](../README.md) pour la documentation complète de l'API.

## 💡 Quick Start

Pour créer votre propre application :

1. **Import** :
   ```go
   import gorulesengine "github.com/deadelus/go-rules-engine/src"
   ```

2. **Engine** :
   ```go
   engine := gorulesengine.NewEngine()
   ```

3. **Règle** :
   ```go
   rule := &gorulesengine.Rule{
       Name:     "my-rule",
       Priority: 100,
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
           Type: "adult",
       },
   }
   engine.AddRule(rule)
   ```

4. **Almanac** :
   ```go
   almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
   almanac.AddFact("age", 25)
   ```

5. **Run** :
   ```go
   results, err := engine.Run(almanac)
   if err != nil {
       log.Fatal(err)
   }
   
   for _, result := range results {
       if result.Result {
           fmt.Printf("✅ %s\n", result.Event.Type)
       }
   }
   ```

## 📝 Structure des exemples

```
examples/
├── README.md           # Ce fichier
├── full-demo.go        # Démo complète
├── basic/              # Exemple basique
│   └── main.go
├── json/               # Chargement JSON
│   └── main.go
├── custom-operator/    # Opérateurs custom
│   └── main.go
└── advanced/           # Features avancées
    └── main.go
```

Consultez chaque exemple pour des cas d'usage spécifiques!

