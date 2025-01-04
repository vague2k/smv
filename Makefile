.PHONY: release
release:
	@if [ -z "$(version)" ]; then \
		echo ""; \
		echo "Error: version is not set. Please specify the version number."; \
		exit 1; \
	fi
	@git tag -a $(version) -m "Release $(version)"
	@git push origin $(version)
