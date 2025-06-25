# Juego de cartas en Go

Este proyecto surgió de la idea de aprender tecnologías nuevas y mejorar las ya existentes. El objetivo principal era desarrollar un juego de cartas en tiempo real mediante websockets (gorilla websockets) para la comunicación bidireccional. El sistema está desarrollado en Go, con una arquitectura hexagonal (ports & adapters) que facilita la separación de responsabilidades, escalabilidad y pruebas. Como base de datos en memoria se utilizó Redis, que permite operaciones de lectura/escritura muy rápidas, lo que permite tener latencias mínimas y permitiría sincronización entre nodos en caso de que se requiriera escalabilidad en el sistema.

![Licencia](https://img.shields.io/badge/Licencia-MIT-blue) ![Estado](https://img.shields.io/badge/Estado-en%20desarrollo-yellow)

---

## 🧠 Tabla de Contenidos

- [Tecnologías](#tecnologías)
- [Características](#características)
- [Instalación](#instalación)

---

## 💻 Tecnologías utilizadas en el Backend

![Go](https://img.shields.io/badge/go-00ADD8.svg?style=for-the-badge&logo=go&logoColor=white) 

![redis](https://img.shields.io/badge/redis-FF4438.svg?style=for-the-badge&logo=redis&logoColor=white)

![docker](https://img.shields.io/badge/docker-2496ED.svg?style=for-the-badge&logo=docker&logoColor=white)

![gin](https://img.shields.io/badge/gin-008ECF.svg?style=for-the-badge&logo=gin&logoColor=white)

---

## ✨ Características

- Escalabilidad
- Mantenibilidad
- Arquitectura hexagonal

---

## 🛠️ Instalación

### Requisitos

- [Just](https://github.com/casey/just) — para automatizar tareas comunes del proyecto
- [Docker](https://www.docker.com/) — para desplegar la app de manera consistente

### 📦 Clonación e instalación

```bash
git clone https://github.com/jorgerr9011/cartas-card-game.git
cd cartas-card-game
just install
just test 
