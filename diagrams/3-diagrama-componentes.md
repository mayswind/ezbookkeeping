# Diagrama de Componentes - Arquitetura de Metas

```mermaid
flowchart TD
    subgraph Cliente ["ezBookkeeping (Cliente)"]
        Frontend["Navegador / PWA"]
        GoalsUI["Módulo de Metas (UI)"]
        Frontend -.->|renderiza| GoalsUI
    end

    subgraph Servidor ["ezBookkeeping (Servidor)"]
        API["API REST"]
        GoalsController["Controlador de Metas"]
        Validator["Validador de Regras"]
        
        API -->|roteia requisições| GoalsController
        GoalsController -->|checa dados| Validator
    end

    subgraph Armazenamento ["Armazenamento"]
        Database[("SQLite / MySQL")]
    end

    GoalsUI -->|requisições HTTP JSON| API
    GoalsController -->|leitura e escrita| Database
