.PHONY: build_start stop login db db_tables clean_all clean_DOCKER_NOOOOOOOOOOOOOOO

# Запуск всего проекта
build_start:
	docker-compose up --build -d

# Остановка и удаление контейнеров
stop:
	docker-compose down

# Зайти внутрь контейнера
login:
	docker exec -it forum_app sh

# ctrl d ыйти типа
db:
	docker exec -it forum_app sqlite3 forum.db -header -column

# Посмотреть только список таблиц
db_tables:
	docker exec -it forum_app sqlite3 forum.db ".tables"

clean_all:
	docker-compose down -v

clean_DOCKER_NOOOOOOOOOOOOOOO:
	docker system prune -a --volumes