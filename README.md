# Projeto de Streaming de Música Pessoal

Este é um projeto de streaming de música para uso pessoal e familiar, com um servidor local e uma interface que permite futuras migrações para armazenamento em nuvem. O objetivo é substituir serviços comerciais como o Spotify, sem fins lucrativos.

## Requisitos Gerais

- **Uso pessoal e familiar**: O sistema deve ser fechado, permitindo acesso apenas a usuários convidados.
- **Autenticação por convites**: Somente usuários com convite podem acessar a plataforma.
- **Interface simples**: O sistema inicial será direto, sem funcionalidades avançadas.
- **Servidor local**: Inicialmente, as músicas serão armazenadas localmente, mas com suporte a futuras migrações para armazenamento em nuvem.
- **Streaming eficiente**: Uso do gRPC para transmissão otimizada de músicas.

## Requisitos do Backend

### 1. Autenticação e Autorização

- Registro de usuários via convite.
- Login e autenticação via e-mail e senha.
- Controle de acesso restrito a usuários registrados.

### 2. Gerenciamento de Músicas

- Upload de arquivos de áudio (MP3, FLAC, etc.).
- Armazenamento de metadados (título, artista, álbum, duração, etc.).
- Estrutura organizada para fácil acesso e busca.
- Suporte a playlists personalizadas.

### 3. Streaming de Música com gRPC

- Serviço gRPC para transmissão de músicas em pacotes pequenos (streaming otimizado).
- Implementação de cache para otimização do desempenho.

### 4. Infraestrutura

- Armazenamento local inicial, com possibilidade de migração para bucket (S3 ou equivalente).
- API RESTful para funcionalidades gerais e gRPC para streaming.
- Arquitetura modular para escalabilidade futura.
- MongoDB como banco de dados.

### 5. Armazenamento

- O sistema de armazenamento pode ser local ou baseado em um bucket S3.
- As músicas e as imagens de álbuns são armazenadas de forma eficiente, permitindo fácil acesso e gerenciamento.
- O armazenamento é gerenciado por uma interface que abstrai a complexidade de interações diretas com o S3 ou o sistema de arquivos local.

### 6. Logs e Monitoramento

- Registros de atividades importantes (login, upload, reprodução de faixas, erros).
- Monitoramento de desempenho do servidor.

## Funcionalidades no Backend

### 1. Autenticação e Autorização

- [x] Registro de usuários via convite.
- [x] Login e autenticação via e-mail e senha.
- [x] Controle de acesso restrito a usuários registrados.

### 2. Gerenciamento de Músicas

- [x] Upload de arquivos de áudio (MP3).
- [ ] Upload de arquivos de letra (texto).
- [ ] Upload de arquivos de cifra (texto).
- [x] Armazenamento de metadados (título, artista, álbum, duração, etc.).
- [x] Busca por música.
- [x] Suporte a playlists personalizadas.

### 3. Streaming de Música com gRPC

- [x] Serviço gRPC para transmissão de músicas em pacotes pequenos (streaming otimizado).
- [ ] Suporte a transmissão de letras e cifras em sincronia com a música.
- [x] Suporte a comunicação bidirecional para comandos de reprodução (play, pause, volume, etc.).
- [ ] Implementação de cache para otimização do desempenho.

### 4. Infraestrutura

- [x] Armazenamento local inicial, com possibilidade de migração para bucket (S3 ou equivalente).
- [x] API RESTful para funcionalidades gerais e gRPC para streaming.

### 5. Logs e Monitoramento

- [ ] Registros de atividades importantes (login, upload, reprodução de faixas, erros).
- [ ] Monitoramento de desempenho do servidor.

## Arquitetura do Sistema

### Diagrama de Contexto

```
+-------------------+
|     Usuário       |
|                   |
| - Interage com o  |
|   aplicativo      |
+-------------------+
          |
          v
+-------------------+
|   Aplicativo      |
|   (Cliente)       |
| - Busca músicas   |
| - Faz upload      |
| - Faz download    |
| - Reproduz música |
+-------------------+
          |
          v
+-------------------+
|   Backend         |
|                   |
| - Servidor gRPC   |
|   - Busca músicas |
|   - Stream de áudio|
|                   |
| - Servidor REST   |
|   - Autenticação  |
|   - Gerenciamento  |
|     de playlists   |
|   - Criação de    |
|     convites      |
+-------------------+
          |
          v
+-------------------+
|   Armazenamento    |
|                   |
| - Local ou S3     |
| - Armazena        |
|   músicas e       |
|   imagens de      |
|   álbuns          |
+-------------------+
          |
          v
+-------------------+
|   Banco de Dados   |
|                   |
| - Armazena        |
|   informações de   |
|   músicas, usuários|
|   e playlists      |
+-------------------+
```

### Diagrama de Contêiner

```
+-------------------+
|   Aplicativo      |
|   (Cliente)       |
|                   |
| - Interface de    |
|   usuário         |
| - Comunicação com  |
|   o backend       |
+-------------------+
          |
          v
+-------------------+
|   Backend         |
|                   |
| - Servidor gRPC   |
|   - Serviço de    |
|     Música        |
|   - Streaming de   |
|     Áudio         |
|                   |
| - Servidor REST   |
|   - Gerenciador de |
|     Autenticação  |
|   - Gerenciador de |
|     Playlists     |
|   - Criação de    |
|     Convites      |
+-------------------+
          |
          v
+-------------------+
|   Armazenamento    |
|                   |
| - Local ou S3     |
| - Armazena        |
|   músicas e       |
|   imagens de      |
|   álbuns          |
+-------------------+
          |
          v
+-------------------+
|   Banco de Dados   |
|                   |
| - MongoDB         |
| - Armazena dados  |
|   de músicas,     |
|   usuários e      |
|   playlists       |
+-------------------+
```

### Diagrama de Arquitetura de Componentes

```
+-------------------+
|   Serviço de Música  |
|                   |
| - uploadMusic()    |
| - downloadMusic()   |
| - searchMusic()     |
| - streamMusic()     |
+-------------------+
          |
          v
+-------------------+
|   Gerenciador de   |
|   Autenticação     |
|                   |
| - login()         |
| - register()      |
| - validateToken() |
+-------------------+
          |
          v
+-------------------+
|   Gerenciador de   |
|   Playlists        |
|                   |
| - createPlaylist() |
| - addMusic()      |
| - removeMusic()   |
| - getPlaylists()  |
+-------------------+
          |
          v
+-------------------+
|   Gerenciador de   |
|   Convites         |
|                   |
| - createInvite()  |
| - validateInvite() |
+-------------------+
          |
          v
+-------------------+
|   Armazenamento    |
|                   |
| - SaveItem()      |
| - GetItem()       |
| - DeleteItem()    |
+-------------------+
```
