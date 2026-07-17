# Diagrama de Componentes - Arquitetura de Metas

```mermaid
componentDiagram
    package "ezBookkeeping (Cliente)" {
        [Navegador / PWA] as Frontend
        [Módulo de Metas (UI)] as GoalsUI
        Frontend ..> GoalsUI : renderiza
    }

    package "ezBookkeeping (Servidor)" {
        [API REST] as API
        [Controlador de Metas] as GoalsController
        [Validador de Regras] as Validator
        
        API --> GoalsController : roteia requisições
        GoalsController --> Validator : checa dados
    }

    database "Armazenamento" {
        [SQLite / MySQL] as Database
    }

    GoalsUI --> API : requisições HTTP (JSON)
    GoalsController --> Database : leitura/escrita
