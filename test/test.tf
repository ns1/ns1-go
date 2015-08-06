provider "nsone" {
}


variable "tld" {
    default = "yelp.tld.example"
}

resource "nsone_datasource" "api" {
    name = "nsone_updater"
    sourcetype = "nsone_v1"
}

resource "nsone_datafeed" "uswest1-prod" {
    name = "uswest1-prod"
    source_id = "${nsone_datasource.api.id}"
    config {
      label = "uswest1-prod"
    }
}

resource "nsone_datafeed" "sfo12-prod" {
    name = "sfo12-prod"
    source_id = "${nsone_datasource.api.id}"
    config {
      label = "sfo12-prod"
    }
}

resource "nsone_datafeed" "dc6-prod" {
    name = "dc6-prod"
    source_id = "${nsone_datasource.api.id}"
    config {
      label = "dc6-prod"
    }
}

resource "nsone_zone" "yelp_tld" {
    zone = "${var.tld}"
    ttl = 60
}

resource "nsone_record" "www" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "www.${var.tld}"
    type = "CNAME"
    answers {
      answer = "www.uswest1.yelp.com"
      meta {
        field = "high_watermark"
        feed = "${nsone_datafeed.uswest1-prod.id}"
      }
      meta {
        field = "low_watermark"
        feed = "${nsone_datafeed.uswest1-prod.id}"
      }
      meta {
        field = "connections"
        feed = "${nsone_datafeed.uswest1-prod.id}"
      }
    }
    answers {
      answer = "www.sfo2.yelp.com"
      meta {
        field = "high_watermark"
        feed = "${nsone_datafeed.sfo12-prod.id}"
      }
      meta {
        field = "low_watermark"
        feed = "${nsone_datafeed.sfo12-prod.id}"
      }
      meta {
        field = "connections"
        feed = "${nsone_datafeed.sfo12-prod.id}"
      }
    }
    answers {
      answer = "www.iad1.yelp.com"
      meta {
        field = "high_watermark"
        feed = "${nsone_datafeed.dc6-prod.id}"
      }
      meta {
        field = "low_watermark"
        feed = "${nsone_datafeed.dc6-prod.id}"
      }
      meta {
        field = "connections"
        feed = "${nsone_datafeed.dc6-prod.id}"
      }
    }
}

    
