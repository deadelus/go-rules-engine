# Go Rules Engine - Architecture

## 1. Vue d'ensemble des composants

```mermaid
graph TB
    subgraph Engine["🎯 Engine (Configuration)"]
        E1[rules: Rule]
        E2[callbacks: map string EventHandler]
        E3[successHandlers: EventHandler]
        E4[failureHandlers: EventHandler]
        E5[eventHandlers: map string EventHandler]
        E6[options: map string interface]
    end

    subgraph Almanac["📋 Almanac (Runtime Ctx)"]
        A1[factMap: map FactID Fact]
        A2[factCache: map string interface]
        A3[events.success: Event]
        A4[events.failure: Event]
        A5[results: RuleResult]
    end

    subgraph Rule["📜 Rule (Business Logic)"]
        R1[Name: string]
        R2[Priority: int]
        R3[Conditions: ConditionSet]
        R4[Event: Event]
        R5[OnSuccess: string]
        R6[OnFailure: string]
    end

    subgraph ConditionSet["🔀 ConditionSet (Logic)"]
        CS1[All: ConditionNode AND]
        CS2[Any: ConditionNode OR]
        CS3[None: ConditionNode NOT]
    end

    subgraph ConditionNode["🔗 ConditionNode"]
        CN1[Condition: Condition]
        CN2[SubSet: ConditionSet]
    end

    subgraph Condition["✓ Condition (Comparison)"]
        C1[Fact: FactID]
        C2[Operator: OperatorType]
        C3[Value: interface]
        C4[Path: string]
        C5[Params: map]
    end

    subgraph Fact["📊 Fact (Data Source)"]
        F1[id: FactID]
        F2[value: interface static]
        F3[calculate: func dynamic]
        F4[options: FactOption]
    end

    subgraph Operator["⚙️ Operator (Comparison)"]
        O1[OperatorEqual]
        O2[OperatorGreaterThan]
        O3[OperatorLessThan]
        O4[OperatorIn]
        O5[OperatorContains]
        O6[OperatorRegex]
        O7[... 11 operators]
    end

    Engine -->|contains| Rule
    Engine -->|Run| Almanac
    Rule -->|has| ConditionSet
    ConditionSet -->|contains| ConditionNode
    ConditionNode -->|leaf| Condition
    ConditionNode -->|nested| ConditionSet
    Condition -->|uses| Operator
    Condition -->|references| Fact
    Almanac -->|stores| Fact
    Almanac -->|caches| Fact
```

## 2. Flux d'exécution

```mermaid
sequenceDiagram
    participant User
    participant Engine
    participant Rule
    participant ConditionSet
    participant Condition
    participant Almanac
    participant Fact
    participant Operator
    participant Handler

    User->>Engine: Run(almanac)
    activate Engine
    
    Engine->>Engine: sortRulesByPriority()
    
    loop For each Rule
        Engine->>Rule: Evaluate Conditions
        activate Rule
        
        Rule->>ConditionSet: Evaluate(almanac)
        activate ConditionSet
        
        loop For each ConditionNode
            ConditionSet->>Condition: Evaluate(almanac)
            activate Condition
            
            Condition->>Almanac: GetFactValue(factID, params, path)
            activate Almanac
            
            alt Fact in cache
                Almanac-->>Condition: cached value
            else Calculate needed
                Almanac->>Fact: Calculate(almanac, params)
                activate Fact
                Fact-->>Almanac: computed value
                deactivate Fact
                Almanac->>Almanac: cache value
                Almanac-->>Condition: computed value
            end
            deactivate Almanac
            
            Condition->>Operator: Evaluate(factValue, compareValue)
            activate Operator
            Operator-->>Condition: bool result
            deactivate Operator
            
            Condition-->>ConditionSet: bool result
            deactivate Condition
        end
        
        ConditionSet->>ConditionSet: Apply All/Any/None logic
        ConditionSet-->>Rule: bool success
        deactivate ConditionSet
        
        Rule-->>Engine: bool success
        deactivate Rule
        
        alt Rule Success
            Engine->>Almanac: AddSuccessEvent(event)
            Engine->>Almanac: AddResult(ruleResult)
            
            opt Rule has OnSuccess callback
                Engine->>Handler: rule.OnSuccess callback
                activate Handler
                Handler-->>Engine: error?
                deactivate Handler
            end
            
            loop Global success handlers
                Engine->>Handler: successHandlers[i]
                activate Handler
                Handler-->>Engine: error?
                deactivate Handler
            end
            
            loop Event-specific handlers
                Engine->>Handler: eventHandlers[event.Type][i]
                activate Handler
                Handler-->>Engine: error?
                deactivate Handler
            end
            
        else Rule Failure
            Engine->>Almanac: AddFailureEvent(event)
            Engine->>Almanac: AddResult(ruleResult)
            
            opt Rule has OnFailure callback
                Engine->>Handler: rule.OnFailure callback
                activate Handler
                Handler-->>Engine: error?
                deactivate Handler
            end
            
            loop Global failure handlers
                Engine->>Handler: failureHandlers[i]
                activate Handler
                Handler-->>Engine: error?
                deactivate Handler
            end
        end
    end
    
    Engine-->>User: []RuleResult, error
    deactivate Engine
```

## 3. Système d'événements

```mermaid
graph TB
    subgraph RuleDefinition["📜 Rule Definition"]
        RD1[Event.Type: string custom]
        RD2[Event.Params: map]
        RD3[OnSuccess: string]
        RD4[OnFailure: string]
    end

    subgraph EngineHandlers["🎯 Engine Handlers"]
        EH1[callbacks: map string EventHandler]
        EH2[successHandlers: EventHandler]
        EH3[failureHandlers: EventHandler]
        EH4[eventHandlers: map string EventHandler]
    end

    subgraph Evaluation["⚖️ Rule Evaluation"]
        EV1{Conditions<br/>Match?}
    end

    subgraph AlmanacStorage["📋 Almanac Storage"]
        AS1[events.success: Event]
        AS2[events.failure: Event]
        AS3[results: RuleResult]
    end

    subgraph HandlerExecution["🔔 Handler Execution Order"]
        HE1["1️⃣ Rule-specific callback<br/>(OnSuccess/OnFailure)"]
        HE2["2️⃣ Global handlers<br/>(successHandlers/failureHandlers)"]
        HE3["3️⃣ Event-type handlers<br/>(eventHandlers[Event.Type])"]
    end

    RuleDefinition --> Evaluation
    
    Evaluation -->|true| AS1
    Evaluation -->|false| AS2
    
    AS1 --> HE1
    AS1 --> HE2
    AS1 --> HE3
    
    AS2 --> HE1
    AS2 --> HE2
    
    EH1 -.->|lookup| HE1
    EH2 -.->|execute all| HE2
    EH3 -.->|execute all| HE2
    EH4 -.->|lookup by type| HE3

    style AS1 fill:#90EE90
    style AS2 fill:#FFB6C1
    style HE1 fill:#FFE4B5
    style HE2 fill:#E0E0E0
    style HE3 fill:#B0C4DE
```

## 4. Structure des conditions (Arbre booléen)

```mermaid
graph TB
    subgraph Example["Exemple: (age >= 18 AND status = 'active') OR (vip = true)"]
        Root[ConditionSet]
        
        Root -->|Any OR| Node1[ConditionNode]
        Root -->|Any OR| Node2[ConditionNode]
        
        Node1 -->|SubSet| CS1[ConditionSet]
        CS1 -->|All AND| CN1[ConditionNode]
        CS1 -->|All AND| CN2[ConditionNode]
        
        CN1 -->|Condition| C1["age >= 18<br/>Fact: 'age'<br/>Operator: OperatorGreaterThanInclusive<br/>Value: 18"]
        CN2 -->|Condition| C2["status = 'active'<br/>Fact: 'status'<br/>Operator: OperatorEqual<br/>Value: 'active'"]
        
        Node2 -->|Condition| C3["vip = true<br/>Fact: 'vip'<br/>Operator: OperatorEqual<br/>Value: true"]
    end

    style Root fill:#FFE4B5
    style CS1 fill:#FFE4B5
    style Node1 fill:#B0C4DE
    style Node2 fill:#B0C4DE
    style CN1 fill:#B0C4DE
    style CN2 fill:#B0C4DE
    style C1 fill:#90EE90
    style C2 fill:#90EE90
    style C3 fill:#90EE90
```

## 5. Types d'opérateurs

```mermaid
graph LR
    subgraph Operators["⚙️ 11 Operators Available"]
        O1[OperatorEqual<br/>'equal']
        O2[OperatorNotEqual<br/>'not_equal']
        O3[OperatorLessThan<br/>'less_than']
        O4[OperatorLessThanInclusive<br/>'less_than_inclusive']
        O5[OperatorGreaterThan<br/>'greater_than']
        O6[OperatorGreaterThanInclusive<br/>'greater_than_inclusive']
        O7[OperatorIn<br/>'in']
        O8[OperatorNotIn<br/>'not_in']
        O9[OperatorContains<br/>'contains']
        O10[OperatorNotContains<br/>'not_contains']
        O11[OperatorRegex<br/>'regex']
    end

    subgraph Registry["Operator Registry"]
        OR[map OperatorType Operator]
    end

    O1 --> OR
    O2 --> OR
    O3 --> OR
    O4 --> OR
    O5 --> OR
    O6 --> OR
    O7 --> OR
    O8 --> OR
    O9 --> OR
    O10 --> OR
    O11 --> OR

    Registry --> Condition[Condition uses<br/>GetOperator opType]

    style O1 fill:#90EE90
    style O2 fill:#90EE90
    style O3 fill:#87CEEB
    style O4 fill:#87CEEB
    style O5 fill:#87CEEB
    style O6 fill:#87CEEB
    style O7 fill:#DDA0DD
    style O8 fill:#DDA0DD
    style O9 fill:#F0E68C
    style O10 fill:#F0E68C
    style O11 fill:#FFB6C1
```

## 6. Builder API (Fluent Interface)

```mermaid
graph LR
    subgraph Builder["🏗️ RuleBuilder (Fluent API)"]
        B1[NewRuleBuilder]
        B2[WithName]
        B3[WithPriority]
        B4[WithConditions]
        B5[WithEvent]
        B6[WithOnSuccess]
        B7[WithOnFailure]
        B8[Build]
    end

    subgraph Helpers["🛠️ Condition Helpers"]
        H1[Equal fact, value]
        H2[GreaterThan fact, value]
        H3[LessThan fact, value]
        H4[In fact, values]
        H5[Contains fact, value]
        H6[Regex fact, pattern]
        H7[... 11 helpers]
    end

    subgraph SetHelpers["🔀 ConditionSet Helpers"]
        S1["All(conditions...) → AND"]
        S2["Any(conditions...) → OR"]
        S3["None(conditions...) → NOT"]
        S4["AllSets(sets...) → Nested AND"]
        S5["AnySets(sets...) → Nested OR"]
        S6["NoneSets(sets...) → Nested NOT"]
    end

    B1 --> B2 --> B3 --> B4 --> B5 --> B6 --> B7 --> B8
    B8 --> Rule[Rule]

    Helpers -.->|used in| B4
    SetHelpers -.->|used in| B4

    style B1 fill:#FFE4B5
    style B8 fill:#90EE90
    style Rule fill:#87CEEB
```

## 7. Fact Types (Statique vs Dynamique)

```mermaid
graph TB
    subgraph FactTypes["📊 Fact Types"]
        FT[Fact]
    end

    FT --> Static["🔒 Static Fact<br/><br/>value: interface{}<br/><br/>Fixed value set at creation"]
    FT --> Dynamic["⚡ Dynamic Fact<br/><br/>calculate: func(almanac, params) interface{}<br/><br/>Computed on-demand"]

    Static --> Cache1["Can be cached<br/>(default: enabled)"]
    Dynamic --> Cache2["Can be cached<br/>(configurable via FactOption)"]

    Cache1 --> Almanac1[Almanac.factCache]
    Cache2 --> Almanac2[Almanac.factCache]

    subgraph Examples["Examples"]
        E1["Static:<br/>age = 25<br/>country = 'FR'<br/>status = 'active'"]
        E2["Dynamic:<br/>isWeekend = func...<br/>userAge = func...<br/>currentTime = func..."]
    end

    Static -.-> E1
    Dynamic -.-> E2

    style Static fill:#90EE90
    style Dynamic fill:#87CEEB
    style Almanac1 fill:#FFE4B5
    style Almanac2 fill:#FFE4B5
```

## 8. Gestion de priorité

```mermaid
graph TB
    subgraph Priority["📊 Rule Priority System"]
        P1[Rule 1: Priority = 100]
        P2[Rule 2: Priority = 50]
        P3[Rule 3: Priority = 75]
    end

    subgraph EngineOptions["⚙️ Engine Options"]
        O1[WithPrioritySorting SortRuleASC]
        O2[WithPrioritySorting SortRuleDESC]
        O3[WithoutPrioritySorting]
    end

    subgraph ExecutionOrder["🔄 Execution Order"]
        E1["Default: DESC<br/>(highest first)"]
        E2["ASC: lowest first"]
        E3["None: insertion order"]
    end

    O1 --> E2
    O2 --> E1
    O3 --> E3

    E1 -->|Sort| Order1["1. Rule 1 (100)<br/>2. Rule 3 (75)<br/>3. Rule 2 (50)"]
    E2 -->|Sort| Order2["1. Rule 2 (50)<br/>2. Rule 3 (75)<br/>3. Rule 1 (100)"]
    E3 -->|No Sort| Order3["1. Rule 1<br/>2. Rule 2<br/>3. Rule 3"]

    style E1 fill:#90EE90
    style E2 fill:#87CEEB
    style E3 fill:#FFE4B5
```

## Légende

- 🎯 Engine: Configuration et orchestration
- 📋 Almanac: Contexte d'exécution runtime
- 📜 Rule: Logique métier
- 🔀 ConditionSet: Groupement logique (AND/OR/NOT)
- ✓ Condition: Comparaison unique
- 📊 Fact: Source de données
- ⚙️ Operator: Logique de comparaison
- 🔔 Handler: Callback d'événement
- 🏗️ Builder: API ergonomique
