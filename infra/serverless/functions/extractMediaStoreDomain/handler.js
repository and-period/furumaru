exports.extractDomainName = async (event) => {
    console.log(event);
    const domainName = event[2].ContainerStatus.Container.Endpoint;
    const domainNameWithoutHttp = domainName.replace('https://', '');
    console.log(domainNameWithoutHttp);
    return domainNameWithoutHttp;
}
