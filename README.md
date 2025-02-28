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
- Suporte a comunicação bidirecional para comandos de reprodução (play, pause, volume, etc.).
- Implementação de cache para otimização do desempenho.

### 4. Infraestrutura

- Armazenamento local inicial, com possibilidade de migração para bucket (S3 ou equivalente).
- API RESTful para funcionalidades gerais e gRPC para streaming.
- Arquitetura modular para escalabilidade futura.

### 5. Logs e Monitoramento

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
- [ ] Estrutura organizada para fácil acesso e busca.
- [ ] Suporte a playlists personalizadas.

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
