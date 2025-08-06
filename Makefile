compile_static:
	g++ -std=c++17 \
	-I./lib/install/include \
	-I./lib/install/include/opencv4 \
	-L./lib/install/lib \
	-lgrpc++ -lgrpc -lgpr -lprotobuf \
	-labsl_log_internal_check_op -labsl_die_if_null -labsl_log_internal_conditions \
	-labsl_log_internal_message -labsl_examine_stack -labsl_log_internal_format \
	-labsl_log_internal_nullguard -labsl_log_internal_structured_proto \
	-labsl_log_internal_proto -labsl_log_internal_log_sink_set -labsl_log_sink \
	-labsl_flags_internal -labsl_flags_marshalling -labsl_flags_reflection \
	-labsl_flags_private_handle_accessor -labsl_flags_commandlineflag \
	-labsl_flags_commandlineflag_internal -labsl_flags_config \
	-labsl_flags_program_name -labsl_log_initialize -labsl_log_internal_globals \
	-labsl_log_globals -labsl_vlog_config_internal -labsl_log_internal_fnmatch \
	-labsl_raw_hash_set -labsl_hash -labsl_city -labsl_low_level_hash \
	-labsl_hashtablez_sampler -labsl_random_distributions -labsl_random_seed_sequences \
	-labsl_random_internal_entropy_pool -labsl_random_internal_randen \
	-labsl_random_internal_randen_hwaes -labsl_random_internal_randen_hwaes_impl \
	-labsl_random_internal_randen_slow -labsl_random_internal_platform \
	-labsl_random_internal_seed_material -labsl_random_seed_gen_exception \
	-labsl_statusor -labsl_status -labsl_cord -labsl_cordz_info -labsl_cord_internal \
	-labsl_cordz_functions -labsl_exponential_biased -labsl_cordz_handle \
	-labsl_crc_cord_state -labsl_crc32c -labsl_crc_internal -labsl_crc_cpu_detect \
	-labsl_leak_check -labsl_strerror -labsl_str_format_internal -labsl_synchronization \
	-labsl_graphcycles_internal -labsl_kernel_timeout_internal -labsl_stacktrace \
	-labsl_symbolize -labsl_debugging_internal -labsl_demangle_internal \
	-labsl_demangle_rust -labsl_decode_rust_punycode -labsl_utf8_for_code_point \
	-labsl_malloc_internal -labsl_tracing_internal -labsl_time -labsl_civil_time \
	-labsl_time_zone -lutf8_validity -lutf8_range -labsl_strings -labsl_strings_internal -labsl_string_view -labsl_int128 -labsl_base -lrt -labsl_spinlock_wait -labsl_throw_delegate -labsl_raw_logging_internal -labsl_log_severity -lopencv_gapi -lopencv_stitching -lopencv_ml -lopencv_objdetect -lopencv_calib3d -lopencv_imgcodecs -lopencv_features2d -lopencv_dnn -lopencv_imgproc -lopencv_core \
	-static-libgcc -static-libstdc++ \
	-pthread -ldl -lz -lssl -lcrypto \
	-o server server.cc service.grpc.pb.cc service.pb.cc

compile_minimal:
	g++ -std=c++17 -pthread -DPROTOBUF_USE_DLLS -DNOMINMAX -I/usr/include/opencv4 -lgrpc++ -lgrpc -lgpr -lprotobuf -labsl_log_internal_check_op -labsl_die_if_null -labsl_log_internal_conditions -labsl_log_internal_message -labsl_examine_stack -labsl_log_internal_format -labsl_log_internal_nullguard -labsl_log_internal_structured_proto -labsl_log_internal_proto -labsl_log_internal_log_sink_set -labsl_log_sink -labsl_flags_internal -labsl_flags_marshalling -labsl_flags_reflection -labsl_flags_private_handle_accessor -labsl_flags_commandlineflag -labsl_flags_commandlineflag_internal -labsl_flags_config -labsl_flags_program_name -labsl_log_initialize -labsl_log_internal_globals -labsl_log_globals -labsl_vlog_config_internal -labsl_log_internal_fnmatch -labsl_raw_hash_set -labsl_hash -labsl_city -labsl_low_level_hash -labsl_hashtablez_sampler -labsl_random_distributions -labsl_random_seed_sequences -labsl_random_internal_entropy_pool -labsl_random_internal_randen -labsl_random_internal_randen_hwaes -labsl_random_internal_randen_hwaes_impl -labsl_random_internal_randen_slow -labsl_random_internal_platform -labsl_random_internal_seed_material -labsl_random_seed_gen_exception -labsl_statusor -labsl_status -labsl_cord -labsl_cordz_info -labsl_cord_internal -labsl_cordz_functions -labsl_exponential_biased -labsl_cordz_handle -labsl_crc_cord_state -labsl_crc32c -labsl_crc_internal -labsl_crc_cpu_detect -labsl_leak_check -labsl_strerror -labsl_str_format_internal -labsl_synchronization -labsl_graphcycles_internal -labsl_kernel_timeout_internal -labsl_stacktrace -labsl_symbolize -labsl_debugging_internal -labsl_demangle_internal -labsl_demangle_rust -labsl_decode_rust_punycode -labsl_utf8_for_code_point -labsl_malloc_internal -labsl_tracing_internal -labsl_time -labsl_civil_time -labsl_time_zone -lutf8_validity -lutf8_range -labsl_strings -labsl_strings_internal -labsl_string_view -labsl_int128 -labsl_base -lrt -labsl_spinlock_wait -labsl_throw_delegate -labsl_raw_logging_internal -labsl_log_severity -lopencv_objdetect -lopencv_calib3d -lopencv_imgcodecs -lopencv_features2d -lopencv_dnn -lopencv_imgproc -lopencv_core -o server server.cc service.grpc.pb.cc service.pb.cc

compile:
	g++ -std=c++17 \
			$(shell pkg-config --cflags --libs grpc++ protobuf opencv4) \
			-o server server.cc service.grpc.pb.cc service.pb.cc

proto-go:
	rm -f client/*.pb.* && \
	protoc --go_out=client --go_opt=paths=source_relative \
    --go-grpc_out=client --go-grpc_opt=paths=source_relative \
		service.proto

proto-cpp:
	rm -f *.pb.* && \
	protoc --grpc_out=. --cpp_out=. service.proto \
		--plugin=protoc-gen-grpc=/usr/bin/grpc_cpp_plugin

proto-cpp-local:
	rm -f *.pb.* && \
	./lib/install/bin/protoc --grpc_out=. --cpp_out=. service.proto \
		--plugin=protoc-gen-grpc=./lib/install/bin/grpc_cpp_plugin
