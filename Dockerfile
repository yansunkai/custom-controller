FROM scratch
ADD custom-controller /custom-controller
ENTRYPOINT ["/custom-controller"]
