#! /bin/bash
if [ -z "${1:-}" ]; then
	echo "Usage: $0 <environment>"
	return
fi

ENV=$1

ENV_FILE=$(ls .env.$ENV* 2>/dev/null | head -n 1)
ENV=$(echo $ENV_FILE | grep -Eo '[^.]+$')

if [ ! -f "$ENV_FILE" ]; then
	echo "Environment file .env.$ENV not found"
	return
fi

set -a
source "$ENV_FILE"
if [ -z "${OLD_PS1:-}" ]; then
	OLD_PS1="$PS1"
fi
PS1="${OLD_PS1:-PS1}(${ENV}) "
set +a
