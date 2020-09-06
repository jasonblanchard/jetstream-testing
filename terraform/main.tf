resource "jetstream_stream" "entries" {
  name     = "ENTRIES"
  subjects = ["ENTRIES.>"]
  storage  = "file"
  max_age  = 60 * 60 * 24 * 365
}

resource "jetstream_consumer" "entries_updated_push" {
  stream_id      = jetstream_stream.entries.id
  durable_name   = "UPDATED_PUSH"
  deliver_all    = true
  filter_subject = "ENTRIES.info.updated"
  sample_freq    = 100
  delivery_subject = "insights.entries.info.updated"
}

resource "jetstream_consumer" "entries_updated_pull" {
  stream_id      = jetstream_stream.entries.id
  durable_name   = "UPDATED_PULL"
  deliver_all    = true
  filter_subject = "ENTRIES.info.updated"
  sample_freq    = 100
}
