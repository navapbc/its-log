from .params import PARAMS


def log_create(cluster, tags, value):
    global PARAMS
    url = f'{PARAMS.base()}/v1/log/create'
    body = {}
    if len(str(cluster)) > 1:
        body["cluster"] = str(cluster)
    if len(str(value)) > 1:
        body["value"] = str(value)
    body["tags"] = tags
    res = PARAMS.session().post(
        url, headers=PARAMS.headers("PUPPERLOG"), json=body)

    return res
