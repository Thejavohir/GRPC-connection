CURRENT_DIR=$(shell pwd)

proto-gen:
	./script/gen-proto.sh ${CURRENT_DIR}
	ls genproto/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp.}"