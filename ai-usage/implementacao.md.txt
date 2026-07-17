Uso de IA na Implementação do MVP
1. Prompts Utilizados

Prompt 1:
"Atue como um desenvolvedor JavaScript. Implemente a lógica para um módulo de Metas Financeiras que permita cadastrar, listar e excluir metas utilizando o localStorage do navegador, com validação dos campos e atualização automática da interface."

Prompt 2:
"Explique o funcionamento do código JavaScript criado, detalhando como utilizar funções, eventos, localStorage e manipulação do DOM para armazenar e exibir as metas financeiras em uma página HTML."

2. Decisões Justificadas e Reflexão Crítica

Por que usamos a IA?

A IA foi utilizada principalmente como apoio no desenvolvimento da lógica em JavaScript, por ser a linguagem em que a equipe possuía menor domínio. O HTML e o CSS foram utilizados apenas para estruturar e estilizar a interface, enquanto a IA auxiliou na implementação das funcionalidades responsáveis pelo cadastro, armazenamento e remoção das metas financeiras.

Reflexão Crítica

A primeira solução sugerida pela IA consistia apenas em uma interface com um botão que exibia uma mensagem de confirmação ao salvar a meta. Após analisar a implementação, percebeu-se que essa abordagem não atendia completamente à história de usuário US01, pois nenhuma informação era realmente armazenada.

Com isso, a solução foi refinada utilizando JavaScript para implementar a lógica da aplicação. Foram desenvolvidas funções para validar os dados informados pelo usuário, armazenar as metas utilizando o localStorage, recuperar automaticamente as metas salvas ao abrir a página e permitir sua exclusão. Também foram adicionadas validações para impedir o cadastro de metas sem nome ou com valores menores ou iguais a zero.

Embora a IA tenha contribuído significativamente para a implementação da lógica em JavaScript, todas as sugestões foram analisadas, adaptadas e testadas pela equipe, garantindo que o código estivesse de acordo com os requisitos definidos para o MVP e com o escopo da disciplina.