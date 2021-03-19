
regen-protos:
	docker run -v `pwd`:/defs namely/protoc-all -i protorepo -d gps/v1beta1/ -l go -o gen

update-protos:
	rm -rf gen/
	docker run -v `pwd`:/defs namely/protoc-all -i protorepo -d gps/v1beta1/ -l go -o gen