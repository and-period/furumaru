exports.extractDomainName = async (event) => {
    console.log(event);
    const domainName = event.ContainerStatus.Container.Endpoint;
    const domainNameWithoutHttp = domainName.replace('https://', '');
    console.log(domainNameWithoutHttp);
    return domainNameWithoutHttp;
}
