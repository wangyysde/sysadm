{{if eq $item.Kind "CHECKBOX" }}
{{$item.Title}}
{{range $ii,$option := $item.ItemData }}
<input type="checkbox"  id="{{$item.ID}}" name="{{$item.Name}}" value="{{- $option.Value -}}" {{if $option.Checked}} checked="checked" {{end}} {{if $item.ActionUri }} onclick='addObjCheckboxClick($item.ActionUri,{{$item.ID}},this.checked)' {{end}}>
{{$option.Text}}
{{end}}
{{end}}



{{if eq $item.Kind "FILE" }}
{{$item.Title}}
<input id="{{$item.ID}}" type="file" name="{{$item.Name}}" {{if gt $item.Size 0}} size="{{$item.Size}}" {{end}}  {{if $item.Disable }} disabled {{end}} value="{{$item.DefaultValue}}" {{if $item.ActionUri }} onblur="addObjvalidFileValue({{$item.ActionUri}},this.value)" {{end}} > {{$item.Note}}
{{end}}

{{if eq $item.Kind "TEXTAREA" }}
{{$item.Title}}
<textarea id="{{$item.ID}}" name="{{$item.Name}}" {{if gt $item.Size 0}} cols="{{$item.Size}}" {{else}} cols="40" {{end}} {{if gt $item.Rows 0}} rows="{{$item.Rows}}" {{else}} rows="5" {{end}} value="{{$item.DefaultValue}}" {{if $item.ActionUri }} onblur="addObjvalidTextareaValue({{$item.ActionUri}},this.value)" {{end}} ></textarea>
{{$item.Note}}
{{end}}

{{if eq $item.Kind "HIDDEN" }}
<input id="{{$item.ID}}" type="hidden" name="{{$item.Name}}"  value="{{$item.DefaultValue}}" >
{{end}}

{{end}}