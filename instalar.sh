#!/bin/bash

set -e

NOME_APP="checklist-rural"
NOME_MENU="Checklist Rural"

PASTA_PROJETO="$(pwd)"
PASTA_INSTALACAO="$HOME/.local/bin/$NOME_APP"
PASTA_MENU="$HOME/.local/share/applications"
EXECUTAVEL="$PASTA_INSTALACAO/$NOME_APP"
SCRIPT_ABRIR="$PASTA_INSTALACAO/abrir-checklist-rural.sh"
ATALHO="$PASTA_MENU/$NOME_APP.desktop"

echo "=========================================="
echo " Instalador - Checklist Rural"
echo "=========================================="
echo

echo "1. Conferindo Go..."
if ! command -v go >/dev/null 2>&1; then
	echo "ERRO: Go não está instalado."
	echo "Instale o Go antes de continuar."
	exit 1
fi

go version
echo

echo "2. Baixando/organizando dependências..."
go mod tidy
echo

echo "3. Compilando o sistema..."
go build -o "$NOME_APP" .
echo

echo "4. Criando pasta de instalação..."
mkdir -p "$PASTA_INSTALACAO"
mkdir -p "$PASTA_MENU"
echo

echo "5. Copiando executável..."
cp "$PASTA_PROJETO/$NOME_APP" "$EXECUTAVEL"
chmod +x "$EXECUTAVEL"
echo

echo "6. Criando script de abertura..."
cat > "$SCRIPT_ABRIR" <<SCRIPT
#!/bin/bash

cd "$PASTA_INSTALACAO" || exit 1

echo "Iniciando Checklist Rural..."
echo "Acesse: http://localhost:8080"
echo
echo "Para fechar o sistema, pressione Ctrl + C."
echo

./$NOME_APP

echo
echo "O sistema foi encerrado."
echo "Pressione Enter para fechar esta janela."
read
SCRIPT

chmod +x "$SCRIPT_ABRIR"
echo

echo "7. Criando atalho no menu..."
cat > "$ATALHO" <<DESKTOP
[Desktop Entry]
Name=$NOME_MENU
Comment=Sistema local para checklist e reuniões de projetos rurais
Exec=x-terminal-emulator -e "$SCRIPT_ABRIR"
Path=$PASTA_INSTALACAO
Icon=applications-office
Terminal=false
Type=Application
Categories=Office;
StartupNotify=false
DESKTOP

chmod +x "$ATALHO"
echo

echo "8. Atualizando menu..."
update-desktop-database "$PASTA_MENU" 2>/dev/null || true
echo

echo "=========================================="
echo " Instalação concluída com sucesso!"
echo "=========================================="
echo
echo "Agora procure no menu por:"
echo "  Checklist Rural"
echo
echo "Ou rode pelo terminal:"
echo "  $EXECUTAVEL"
echo
