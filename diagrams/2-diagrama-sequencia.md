# Diagrama de Sequência - Criar Nova Meta (US01)

```mermaid
sequenceDiagram
    actor Estudante
    participant UI as Interface Web (Frontend)
    participant API as Backend (Servidor)
    participant DB as Banco de Dados

    Estudante->>UI: Acessa aba "Metas" e clica em "Nova Meta"
    UI-->>Estudante: Exibe formulário
    Estudante->>UI: Preenche Nome (Ex: Congresso ENEJ) e Valor (R$ 800)
    UI->>API: POST /api/goals {name, target_amount}
    
    alt Dados Válidos
        API->>DB: INSERT INTO goals (name, target)
        DB-->>API: Confirmação de salvamento
        API-->>UI: Retorna Status 201 (Created)
        UI-->>Estudante: Exibe mensagem "Meta criada com sucesso" e atualiza lista
    else Dados Inválidos (Ex: Valor negativo)
        API-->>UI: Retorna Status 400 (Bad Request)
        UI-->>Estudante: Exibe mensagem de erro "Valor deve ser maior que zero"
    end
