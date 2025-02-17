generate_ddl:
	@go run tools/db_gen/ddl/main.go

generate_er_puml:
	@go run tools/db_gen/er/main.go

# 変数定義
ARCH=$(shell uname -m)
ROOT_DIR=$(PWD)
PLANTUML_JAR=plantuml.jar
PLANTUML_CMD=java -DPLANTUML_LIMIT_SIZE=16384 -jar $(PLANTUML_JAR)
ER_OUTPUT_DIR=$(ROOT_DIR)/db/er/image
DOCKER_IMAGE=er_image_generator

# 出力ファイル
USER_ER_OUTPUT=$(ER_OUTPUT_DIR)/er_user_db_gen.svg
MASTER_ER_OUTPUT=$(ER_OUTPUT_DIR)/er_master_db_gen.svg

# Docker イメージの設定（アーキテクチャごとに異なる）
ifeq ($(ARCH),x86_64)
    BASE_IMAGE=openjdk:19-jdk-alpine3.16
else ifeq ($(ARCH),arm64)
    BASE_IMAGE=arm64v8/openjdk:8-jre-alpine
else
    $(error "Unsupported architecture: $(ARCH)")
endif

# Docker イメージをビルド
.PHONY: build_er_docker
build_er_docker:
	docker build -t $(DOCKER_IMAGE) --build-arg IMAGE=$(BASE_IMAGE) -f docker/puml/Dockerfile docker/puml/

# ユーザー用 ER図の SVG を生成
.PHONY: generate_user_er_svg
generate_user_er_svg: build_er_docker
	@mkdir -p $(ER_OUTPUT_DIR)
	docker run --rm -v $(ROOT_DIR)/db/er:/er $(DOCKER_IMAGE) \
		$(PLANTUML_CMD) -charset UTF-8 -nometadata -nbthread auto -progress -t"svg" -o "../image" "$(ROOT_DIR)/db/er/er_user_db_gen.puml"

# マスター用 ER図の SVG を生成
.PHONY: generate_master_er_svg
generate_master_er_svg: build_er_docker
	@mkdir -p $(ER_OUTPUT_DIR)
	docker run --rm -v $(ROOT_DIR)/db/er:/er $(DOCKER_IMAGE) \
		$(PLANTUML_CMD) -charset UTF-8 -nometadata -nbthread auto -progress -t"svg" -o "../image" "$(ROOT_DIR)/db/er/er_master_db_gen.puml"

# すべての ER図を生成
.PHONY: generate_er_svg
generate_er_svg: generate_user_er_svg generate_master_er_svg

# クリーンアップ
.PHONY: clean
clean:
	rm -rf $(ER_OUTPUT_DIR)
	docker rmi $(DOCKER_IMAGE)
