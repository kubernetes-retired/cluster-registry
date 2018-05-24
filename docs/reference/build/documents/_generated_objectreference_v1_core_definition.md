## ObjectReference v1

Group        | Version     | Kind
------------ | ---------- | -----------
`core` | `v1` | `ObjectReference`



ObjectReference contains enough information to let you inspect or modify the referred object.

<aside class="notice">
Appears In:

<ul> 
<li><a href="#authinfo-v1alpha1">AuthInfo v1alpha1</a></li>
</ul></aside>

Field        | Description
------------ | -----------
`apiVersion`<br /> *string*    | API version of the referent.
`fieldPath`<br /> *string*    | If referring to a piece of an object instead of an entire object, this string should contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2]. For example, if the object reference is to a container within a pod, this would take on a value like: &#34;spec.containers{name}&#34; (where &#34;name&#34; refers to the name of the container that triggered the event) or if no container name is specified &#34;spec.containers[2]&#34; (container with index 2 in this pod). This syntax is chosen only to have some well-defined way of referencing a part of an object.
`kind`<br /> *string*    | Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
`name`<br /> *string*    | Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
`namespace`<br /> *string*    | Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
`resourceVersion`<br /> *string*    | Specific resourceVersion to which this reference is made, if any. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#concurrency-control-and-consistency
`uid`<br /> *string*    | UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids

