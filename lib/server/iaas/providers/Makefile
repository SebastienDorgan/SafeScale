GO?=go

.PHONY:	clean generate vet

generate:
	@(cd cloudferro && $(MAKE) $(@))
	@(cd flexibleengine && $(MAKE) $(@))
	@(cd openstack && $(MAKE) $(@))
	@(cd opentelekom && $(MAKE) $(@))
	@(cd ovh && $(MAKE) $(@))
	@(cd local && $(MAKE) $(@))
	@(cd gcp && $(MAKE) $(@))
	@(cd aws && $(MAKE) $(@))

vet:
	@$(GO) vet ./...

clean:
	@(cd cloudferro && $(MAKE) $(@))
	@(cd flexibleengine && $(MAKE) $(@))
	@(cd openstack && $(MAKE) $(@))
	@(cd opentelekom && $(MAKE) $(@))
	@(cd ovh && $(MAKE) $(@))
	@(cd local && $(MAKE) $(@))
	@(cd gcp && $(MAKE) $(@))
	@(cd aws && $(MAKE) $(@))
	@$(RM) ./mocks/*.go || true
