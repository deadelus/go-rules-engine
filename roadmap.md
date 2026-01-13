Parfait ! Voici une roadmap détaillée et adaptée pour un agent IA comme Roo Code, avec des instructions très précises pour chaque étape :

## Roadmap pour Agent IA - Adaptation json-rules-engine en Go

### 🎯 Phase 1 : Setup et structures de base

**ÉTAPE 1.1 : Initialisation du projet**
```
Objectif : Créer la structure du projet Go

Actions à réaliser :
1. Créer le répertoire du projet : go-rules-engine/
2. Initialiser le module Go : go mod init github.com/deadelus/go-rules-engine
3. Créer les fichiers suivants (vides pour l'instant) :
   - engine.go
   - rule.go
   - condition.go
   - operator.go
   - fact.go
   - event.go
   - almanac.go
   - types.go
   - errors.go

Validation : Le projet compile avec "go build"
```

**ÉTAPE 1.2 : Définir les types de base dans types.go**
```
Objectif : Créer les enums et types fondamentaux

Fichier : types.go
Créer :
1. Type OperatorType (string) avec les constantes :
   - OperatorEqual = "equal"
   - OperatorNotEqual = "notEqual"
   - OperatorLessThan = "lessThan"
   - OperatorLessThanInclusive = "lessThanInclusive"
   - OperatorGreaterThan = "greaterThan"
   - OperatorGreaterThanInclusive = "greaterThanInclusive"
   - OperatorIn = "in"
   - OperatorNotIn = "notIn"
   - OperatorContains = "contains"
   - OperatorDoesNotContain = "doesNotContain"

2. Type ConditionType (string) avec :
   - ConditionAll = "all"
   - ConditionAny = "any"

3. Type Fact = map[string]interface{} (alias)

Validation : Le fichier compile sans erreur
```

**ÉTAPE 1.3 : Créer la structure Condition dans condition.go**
```
Objectif : Définir comment représenter une condition

Fichier : condition.go
Créer la struct Condition avec les champs :
- Fact string (le chemin vers la propriété, ex: "user.age")
- Operator OperatorType (l'opérateur à utiliser)
- Value interface{} (la valeur à comparer)
- Path string (optionnel, pour traverser les objets)
- Params map[string]interface{} (optionnel, paramètres additionnels)

Ajouter les tags JSON pour la sérialisation :
`json:"fact"`, `json:"operator"`, etc.

Validation : Pouvoir unmarshaler ce JSON :
{
  "fact": "age",
  "operator": "greaterThan",
  "value": 18
}
```

**ÉTAPE 1.4 : Créer la structure ConditionSet dans condition.go**
```
Objectif : Gérer les conditions composées (AND/OR)

Fichier : condition.go
Créer la struct ConditionSet avec :
- Type ConditionType (all ou any)
- Conditions []ConditionProperties
- All []ConditionProperties (optionnel)
- Any []ConditionProperties (optionnel)

Créer le type ConditionProperties qui peut être soit :
- Une Condition simple
- Un ConditionSet imbriqué

Astuce : Utiliser interface{} ou créer un type union avec json.RawMessage

Validation : Pouvoir unmarshaler ce JSON :
{
  "all": [
    {"fact": "age", "operator": "greaterThan", "value": 18},
    {"fact": "country", "operator": "equal", "value": "FR"}
  ]
}
```

**ÉTAPE 1.5 : Créer la structure Rule dans rule.go**
```
Objectif : Représenter une règle complète

Fichier : rule.go
Créer la struct Rule avec :
- Name string (nom de la règle, optionnel)
- Priority int (priorité d'exécution, défaut: 1)
- Conditions ConditionSet (les conditions à évaluer)
- Event Event (l'événement déclenché si succès)
- OnSuccess *string (optionnel, callback name)
- OnFailure *string (optionnel, callback name)

Ajouter les tags JSON appropriés

Validation : Pouvoir unmarshaler une règle JSON complète
```

### 🔧 Phase 2 : Système d'opérateurs

**ÉTAPE 2.1 : Interface Operator dans operator.go**
```
Objectif : Définir le contrat pour tous les opérateurs

Fichier : operator.go
Créer l'interface Operator avec la méthode :
- Evaluate(factValue interface{}, compareValue interface{}) (bool, error)

Cette méthode prend la valeur du fait et la valeur à comparer,
et retourne true si la condition est satisfaite.

Validation : L'interface compile
```

**ÉTAPE 2.2 : Implémenter les opérateurs de comparaison**
```
Objectif : Créer les opérateurs ==, !=, <, >, <=, >=

Fichier : operator.go
Pour chaque opérateur, créer une struct qui implémente Operator :

1. equalOperator struct
   - Implement Evaluate() qui compare l'égalité
   - Gérer les types : int, float64, string, bool
   - Utiliser reflect.DeepEqual pour les autres types

2. notEqualOperator struct
   - Inverse de equalOperator

3. lessThanOperator struct
   - Comparer des nombres (int, float64)
   - Retourner une erreur si types incompatibles

4. greaterThanOperator struct
   - Similaire à lessThan mais inverse

5. lessThanInclusiveOperator (<=)
6. greaterThanInclusiveOperator (>=)

Astuce : Créer une fonction helper compareNumbers(a, b interface{}) pour normaliser les types numériques

Validation : Écrire des tests unitaires pour chaque opérateur
```

**ÉTAPE 2.3 : Implémenter les opérateurs de collection**
```
Objectif : Créer in, notIn, contains, doesNotContain

Fichier : operator.go

1. inOperator struct
   - Evaluate() vérifie si factValue est dans le slice compareValue
   - compareValue doit être un []interface{}
   - Supporter les types de base

2. notInOperator struct
   - Inverse de inOperator

3. containsOperator struct
   - factValue doit être un slice ou string
   - Vérifie si compareValue est dedans
   - Pour string : utiliser strings.Contains()
   - Pour slice : itérer et comparer

4. doesNotContainOperator struct
   - Inverse de containsOperator

Validation : Tests avec différents types de données
```

**ÉTAPE 2.4 : Créer la registry d'opérateurs**
```
Objectif : Centraliser tous les opérateurs disponibles

Fichier : operator.go
Créer :
1. Une map global operatorRegistry = map[OperatorType]Operator{}

2. Une fonction init() qui enregistre tous les opérateurs :
   operatorRegistry[OperatorEqual] = &equalOperator{}
   operatorRegistry[OperatorNotEqual] = &notEqualOperator{}
   ... etc

3. Une fonction GetOperator(op OperatorType) (Operator, error)
   qui retourne l'opérateur depuis la registry

4. Une fonction RegisterOperator(name OperatorType, op Operator)
   pour permettre d'ajouter des opérateurs custom

Validation : Pouvoir récupérer n'importe quel opérateur par son nom
```

### 📊 Phase 3 : Almanac (gestion des faits)

**ÉTAPE 3.1 : Créer la structure Almanac dans almanac.go**
```
Objectif : Stocker et récupérer les faits efficacement

Fichier : almanac.go
Créer la struct Almanac avec :
- facts map[string]interface{} (stockage des faits)
- factResults map[string]interface{} (cache des résultats calculés)
- mutex sync.RWMutex (pour la concurrence)

Créer les méthodes :
1. NewAlmanac(facts Fact) *Almanac
   - Constructeur qui initialise les maps

2. AddFact(path string, value interface{})
   - Ajoute un fait au stockage

3. GetFactValue(path string) (interface{}, error)
   - Récupère la valeur d'un fait par son path
   - Gérer les paths simples ("age") et imbriqués ("user.address.city")

Validation : Tests avec des faits simples et imbriqués
```

**ÉTAPE 3.2 : Implémenter le path traversal**
```
Objectif : Naviguer dans les objets imbriqués

Fichier : almanac.go
Créer la fonction helper :
- traversePath(data interface{}, path string) (interface{}, error)

Logic :
1. Split le path par "." : []string{"user", "address", "city"}
2. Pour chaque segment :
   - Vérifier que data est une map[string]interface{}
   - Extraire la clé suivante
   - Continuer avec la valeur extraite
3. Retourner la valeur finale

Gérer les cas d'erreur :
- Path inexistant
- Type incompatible (pas une map)
- Nil values

Intégrer cette fonction dans GetFactValue()

Validation : Tests avec des structures profondément imbriquées
```

**ÉTAPE 3.3 : Ajouter le support des faits dynamiques**
```
Objectif : Permettre des faits calculés à la volée

Fichier : almanac.go
Ajouter à Almanac :
- factFunctions map[string]FactFunction
- type FactFunction = func(*Almanac) (interface{}, error)

Créer les méthodes :
1. AddFactFunction(name string, fn FactFunction)
   - Enregistre une fonction qui calcule un fait

2. Modifier GetFactValue() pour :
   - Vérifier d'abord dans le cache factResults
   - Si pas trouvé, chercher dans facts
   - Si pas trouvé, chercher dans factFunctions et exécuter
   - Mettre en cache le résultat

Validation : Créer un fait dynamique "currentHour" qui retourne time.Now().Hour()
```

### ⚙️ Phase 4 : Évaluation des conditions

**ÉTAPE 4.1 : Créer l'évaluateur de condition simple**
```
Objectif : Évaluer une condition unique

Fichier : condition.go
Créer la méthode :
- (c *Condition) Evaluate(almanac *Almanac) (bool, error)

Logic :
1. Récupérer la valeur du fait : almanac.GetFactValue(c.Fact)
2. Récupérer l'opérateur : GetOperator(c.Operator)
3. Appeler operator.Evaluate(factValue, c.Value)
4. Retourner le résultat

Gérer les erreurs :
- Fait non trouvé
- Opérateur invalide
- Erreur d'évaluation

Validation : Tests unitaires avec différentes conditions
```

**ÉTAPE 4.2 : Créer l'évaluateur de ConditionSet**
```
Objectif : Évaluer les conditions composées (all/any)

Fichier : condition.go
Créer la méthode :
- (cs *ConditionSet) Evaluate(almanac *Almanac) (bool, error)

Logic pour "all" :
1. Parcourir toutes les conditions
2. Si une condition retourne false, retourner false immédiatement
3. Si toutes retournent true, retourner true

Logic pour "any" :
1. Parcourir toutes les conditions
2. Si une condition retourne true, retourner true immédiatement
3. Si toutes retournent false, retourner false

Gérer les conditions imbriquées (récursion)

Validation : Tests avec conditions composées et imbriquées
```

**ÉTAPE 4.3 : Ajouter le support des paramètres**
```
Objectif : Permettre des conditions paramétrables

Fichier : condition.go
Modifier Condition.Evaluate() pour :
1. Vérifier si c.Params contient des clés
2. Injecter les params dans l'almanac temporairement
3. Évaluer la condition
4. Nettoyer les params après évaluation

Exemple d'utilisation :
{
  "fact": "temperature",
  "operator": "greaterThan",
  "value": {"fact": "threshold"},
  "params": {
    "threshold": 25
  }
}

Validation : Tests avec paramètres dynamiques
```

### 🎯 Phase 5 : Système d'événements

**ÉTAPE 5.1 : Créer la structure Event dans event.go**
```
Objectif : Représenter un événement déclenché

Fichier : event.go
Créer la struct Event avec :
- Type string (nom de l'événement, ex: "user-adult")
- Params map[string]interface{} (données additionnelles)

Créer la struct EventResult avec :
- Event Event
- Rule *Rule (la règle qui a déclenché)
- Almanac *Almanac (accès aux faits)
- Result bool (succès ou échec)

Ajouter les tags JSON

Validation : Structure compile correctement
```

**ÉTAPE 5.2 : Créer le gestionnaire d'événements**
```
Objectif : Permettre d'enregistrer des handlers

Fichier : event.go
Créer :
1. type EventHandler = func(EventResult) error

2. Struct EventEmitter avec :
   - handlers map[string][]EventHandler
   - mutex sync.RWMutex

3. Méthodes de EventEmitter :
   - On(eventType string, handler EventHandler)
     → Enregistre un handler
   
   - Emit(result EventResult) error
     → Exécute tous les handlers pour cet événement
     → Retourne la première erreur rencontrée

   - Off(eventType string)
     → Supprime tous les handlers d'un type

Validation : Tests d'enregistrement et d'émission d'événements
```

### 🚀 Phase 6 : Moteur principal

**ÉTAPE 6.1 : Créer la structure Engine dans engine.go**
```
Objectif : Le cœur du moteur de règles

Fichier : engine.go
Créer la struct Engine avec :
- rules []*Rule (liste des règles)
- emitter *EventEmitter (pour les événements)
- operators map[OperatorType]Operator (registry locale)
- allowUndefinedFacts bool (option pour gérer faits manquants)

Créer le constructeur :
- NewEngine(options ...EngineOption) *Engine
  → Initialise l'engine avec des options fonctionnelles

Créer les options :
- type EngineOption = func(*Engine)
- WithAllowUndefinedFacts(allow bool) EngineOption
- WithOperator(op OperatorType, operator Operator) EngineOption

Validation : Pouvoir créer un engine vide
```

**ÉTAPE 6.2 : Ajouter/supprimer des règles**
```
Objectif : Gérer la collection de règles

Fichier : engine.go
Créer les méthodes :

1. (e *Engine) AddRule(rule *Rule) error
   - Valide la règle
   - L'ajoute à e.rules
   - Trie par priorité (plus haute en premier)

2. (e *Engine) RemoveRule(ruleName string) error
   - Cherche la règle par nom
   - La supprime de la liste

3. (e *Engine) AddRuleFromJSON(jsonData []byte) error
   - Unmarshal le JSON vers Rule
   - Appelle AddRule()

4. Fonction helper sortRulesByPriority(rules []*Rule)
   - Utilise sort.Slice()
   - Priorité décroissante

Validation : Tests d'ajout/suppression de règles
```

**ÉTAPE 6.3 : Implémenter Engine.Run()**
```
Objectif : Exécuter toutes les règles contre des faits

Fichier : engine.go
Créer la méthode :
- (e *Engine) Run(facts Fact) ([]EventResult, error)

Logic :
1. Créer un Almanac avec les facts
2. Créer un slice results []EventResult
3. Pour chaque règle (dans l'ordre de priorité) :
   a. Évaluer rule.Conditions.Evaluate(almanac)
   b. Si true :
      - Créer un EventResult avec success=true
      - Émettre l'événement : e.emitter.Emit(result)
      - Ajouter à results
   c. Si false :
      - Créer un EventResult avec success=false
      - Ajouter à results si OnFailure défini
4. Retourner results et nil

Gérer les erreurs d'évaluation

Validation : Tests end-to-end avec plusieurs règles
```

**ÉTAPE 6.4 : Ajouter les callbacks sur succès/échec**
```
Objectif : Exécuter des actions après évaluation

Fichier : engine.go
Modifier Engine pour ajouter :
- successCallbacks map[string]EventHandler
- failureCallbacks map[string]EventHandler

Créer les méthodes :
1. (e *Engine) OnSuccess(name string, handler EventHandler)
2. (e *Engine) OnFailure(name string, handler EventHandler)

Modifier Run() pour :
- Après évaluation d'une règle, si OnSuccess défini :
  → Chercher dans successCallbacks et exécuter
- Pareil pour OnFailure

Validation : Tests avec callbacks
```

### 🔍 Phase 7 : Validation et erreurs

**ÉTAPE 7.1 : Créer les erreurs personnalisées dans errors.go**
```
Objectif : Avoir des erreurs claires et typées

Fichier : errors.go
Créer les types d'erreur :

1. type ErrInvalidOperator struct { Operator string }
   - func (e ErrInvalidOperator) Error() string

2. type ErrFactNotFound struct { Path string }
   - func (e ErrFactNotFound) Error() string

3. type ErrInvalidCondition struct { Reason string }
   - func (e ErrInvalidCondition) Error() string

4. type ErrInvalidRule struct { RuleName string; Reason string }
   - func (e ErrInvalidRule) Error() string

5. type ErrTypeConversion struct { From, To string }
   - func (e ErrTypeConversion) Error() string

Utiliser ces erreurs dans tout le code

Validation : Tests de gestion d'erreurs
```

**ÉTAPE 7.2 : Créer le validateur de règles**
```
Objectif : Valider les règles avant exécution

Fichier : rule.go
Créer la méthode :
- (r *Rule) Validate() error

Vérifications :
1. Conditions ne doit pas être vide
2. Event.Type ne doit pas être vide
3. Valider récursivement toutes les conditions :
   - Fact non vide
   - Operator valide (existe dans registry)
   - Value non nil (sauf pour certains opérateurs)
4. Priority doit être >= 1

Retourner ErrInvalidRule avec détails

Appeler Validate() dans Engine.AddRule()

Validation : Tests avec règles invalides
```

### 📚 Phase 8 : API ergonomique et builders

**ÉTAPE 8.1 : Créer des builders fluent**
```
Objectif : Faciliter la création de règles en code

Fichier : builder.go
Créer RuleBuilder :

type RuleBuilder struct {
    rule *Rule
}

Méthodes chainables :
1. NewRuleBuilder() *RuleBuilder
2. (rb *RuleBuilder) WithName(name string) *RuleBuilder
3. (rb *RuleBuilder) WithPriority(p int) *RuleBuilder
4. (rb *RuleBuilder) WithCondition(cond *Condition) *RuleBuilder
5. (rb *RuleBuilder) WithEvent(eventType string, params map[string]interface{}) *RuleBuilder
6. (rb *RuleBuilder) Build() (*Rule, error)

Exemple d'utilisation :
rule, _ := NewRuleBuilder().
    WithName("adult-user").
    WithPriority(10).
    WithCondition(&Condition{...}).
    WithEvent("user-is-adult", nil).
    Build()

Validation : Tests de construction de règles
```

**ÉTAPE 8.2 : Créer des helpers de conditions**
```
Objectif : Simplifier la création de conditions

Fichier : builder.go
Créer des fonctions :

1. Equal(fact string, value interface{}) *Condition
2. NotEqual(fact string, value interface{}) *Condition
3. GreaterThan(fact string, value interface{}) *Condition
4. LessThan(fact string, value interface{}) *Condition
5. In(fact string, values []interface{}) *Condition
6. Contains(fact string, value interface{}) *Condition
7. All(conditions ...*Condition) *ConditionSet
8. Any(conditions ...*Condition) *ConditionSet

Exemple d'utilisation :
cond := All(
    GreaterThan("age", 18),
    Equal("country", "FR"),
)

Validation : Tests de création de conditions
```

### 📖 Phase 9 : Documentation et exemples

**ÉTAPE 9.1 : Documenter le code**
```
Objectif : Ajouter les commentaires GoDoc

Pour chaque fichier :
1. Ajouter un commentaire de package en haut
2. Documenter chaque struct exportée
3. Documenter chaque méthode/fonction exportée
4. Ajouter des exemples en commentaire

Format GoDoc :
// Engine représente le moteur de règles principal.
// Il contient une collection de règles et les exécute
// contre un ensemble de faits fournis.
type Engine struct { ... }

Validation : go doc devrait afficher la documentation
```

**ÉTAPE 9.2 : Créer des exemples dans examples/**
```
Objectif : Fournir des cas d'usage concrets

Créer examples/basic/main.go :
- Exemple simple avec une règle
- Vérifier l'âge d'un utilisateur

Créer examples/json/main.go :
- Charger des règles depuis JSON
- Exécuter et afficher les résultats

Créer examples/advanced/main.go :
- Règles multiples avec priorités
- Conditions imbriquées
- Handlers d'événements
- Faits dynamiques

Créer examples/custom-operator/main.go :
- Créer un opérateur personnalisé
- L'enregistrer dans l'engine

Chaque exemple doit être exécutable avec "go run"

Validation : Tous les exemples compilent et s'exécutent
```

**ÉTAPE 9.3 : Écrire le README.md**
```
Objectif : Documentation complète pour les utilisateurs

Sections du README :
1. Titre et description
2. Installation (go get)
3. Quick Start (exemple minimal)
4. Concepts clés (Rules, Conditions, Facts, Events)
5. Usage détaillé :
   - Créer un engine
   - Ajouter des règles
   - Exécuter le moteur
   - Gérer les événements
6. Opérateurs disponibles (tableau)
7. API Reference (lien vers GoDoc)
8. Exemples avancés
9. Différences avec json-rules-engine JS
10. Contribution et licence

Validation : README clair et complet
```

### ✅ Phase 10 : Tests et qualité

**ÉTAPE 10.1 : Tests unitaires complets**
```
Objectif : Couvrir tout le code avec des tests

Pour chaque fichier .go, créer un fichier _test.go :
- operator_test.go : tester chaque opérateur
- condition_test.go : tester l'évaluation
- rule_test.go : tester la validation
- almanac_test.go : tester le path traversal
- engine_test.go : tester l'exécution complète

Utiliser table-driven tests :
func TestEqualOperator(t *testing.T) {
    tests := []struct{
        name string
        factValue interface{}
        compareValue interface{}
        expected bool
    }{
        {"int equal", 5, 5, true},
        {"int not equal", 5, 10, false},
        ...
    }
    for _, tt := range tests { ... }
}

Objectif : > 80% de couverture

Validation : go test -cover ./...
```

**ÉTAPE 10.2 : Tests d'intégration**
```
Objectif : Tester des scénarios complets

Créer engine_integration_test.go :

1. Test avec règles multiples et priorités
2. Test avec conditions complexes imbriquées
3. Test avec faits dynamiques
4. Test avec événements et handlers
5. Test avec JSON complet (unmarshaling + exécution)
6. Test de performance avec beaucoup de règles

Validation : Tous les tests passent
```

**ÉTAPE 10.3 : Benchmarks**
```
Objectif : Mesurer les performances

Créer des benchmarks dans *_test.go :

1. BenchmarkSimpleCondition
2. BenchmarkComplexConditions
3. BenchmarkEngineRun (1 règle)
4. BenchmarkEngineRun (100 règles)
5. BenchmarkEngineRun (1000 règles)
6. BenchmarkPathTraversal

Utiliser b.ResetTimer() et b.ReportAllocs()

Validation : go test -bench=. -benchmem
```

### 🎁 Phase 11 : Finalisation

**ÉTAPE 11.1 : Ajouter les fichiers du projet**
```
Objectif : Compléter le repository

Créer :
1. LICENSE (MIT ou autre)
2. .gitignore (fichiers Go standard)
3. CHANGELOG.md (v0.1.0 initial release)
4. CONTRIBUTING.md (guidelines de contribution)
5. go.mod et go.sum à jour

Validation : Structure complète du projet
```

**ÉTAPE 11.2 : Release v1.0.0**
```
Objectif : Préparer la première version stable

Checklist :
☐ Tous les tests passent
☐ Documentation complète
☐ Exemples fonctionnels
☐ README à jour
☐ Version dans go.mod
☐ Tag git v1.0.0
☐ GitHub release avec notes

Commandes :
git tag v1.0.0
git push origin v1.0.0

Validation : Le package est utilisable via go get
```

---

## 🎯 Ordre d'exécution recommandé pour l'agent IA

Voici l'ordre optimal que l'agent devrait suivre :

1. **Semaine 1** : Phases 1-2 (Setup + Opérateurs)
2. **Semaine 2** : Phases 3-4 (Almanac + Évaluation)
3. **Semaine 3** : Phases 5-6 (Events + Engine)
4. **Semaine 4** : Phases 7-8 (Validation + API)
5. **Semaine 5** : Phases 9-10 (Documentation + Tests)
6. **Semaine 6** : Phase 11 (Finalisation)

## 📋 Checklist pour l'agent après chaque phase

Après chaque phase, l'agent doit vérifier :
- ✅ Le code compile sans erreur
- ✅ Les tests unitaires passent
- ✅ La documentation est à jour
- ✅ Pas de TODO ou FIXME critiques
- ✅ Le code suit les conventions Go (gofmt, golint)