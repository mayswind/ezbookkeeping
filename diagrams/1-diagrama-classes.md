# Diagrama de Classes - MVP (Metas Financeiras)

```mermaid
classDiagram
    class User {
        +int id
        +String username
        +String password
        +login()
        +logout()
    }

    class Account {
        +int id
        +String name
        +float balance
        +addFunds(amount)
        +deductFunds(amount)
    }

    class Goal {
        +int id
        +String name
        +float targetAmount
        +float currentAmount
        +createGoal()
        +allocateFunds(amount)
        +getProgress()
    }

    User "1" -- "*" Account : possui
    User "1" -- "*" Goal : define
    Account "1" -- "*" Goal : financia
