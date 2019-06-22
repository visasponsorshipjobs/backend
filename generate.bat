protoc job\jobpb\job.proto --go_out=plugins=grpc:.

protoc backend/job/jobpb/job.proto --js_out=import_style=commonjs,binary:web/web/web-client/src/