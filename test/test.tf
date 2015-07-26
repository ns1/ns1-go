provider "nsone" {
}

resource "nsone_zone" "yelp_co_uk" {
    zone = "yelp.co.uk"
    primary = "199.255.189.133"
}

variable "tld" {
    default = "yelp.tld.example"
}

resource "nsone_datasource" "api" {
    name = "terraform test"
    sourcetype = "nsone_v1"
}

resource "nsone_datafeed" "exampledc1" {
    name = "exampledc3"
    source_id = "${nsone_datasource.api.id}"
    config {
      label = "exampledc1"
    }
}

resource "nsone_datafeed" "exampledc3" {
    name = "exampledc3"
    source_id = "${nsone_datasource.api.id}"
    config {
      label = "exampledc3"
    }
}

resource "nsone_zone" "yelp_tld" {
    zone = "${var.tld}"
    ttl = 60
}

resource "nsone_zone" "linktest" {
    zone = "link.tld.example"
    link = "${nsone_zone.yelp_tld.zone}"
}

resource "nsone_record" "test" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "test.${var.tld}"
    type = "A"
    link = "www.${var.tld}"
}

resource "nsone_record" "star" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "*.${var.tld}"
    type = "A"
    answers { 
      answer = "198.51.132.28"
    }
    answers {
      answer = "198.51.132.228"
    }
}

resource "nsone_record" "www" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "www.${var.tld}"
    type = "A"
    answers {
      answer = "198.51.132.28"
      meta {
        field = "up"
        feed = "${nsone_datafeed.exampledc1.id}"
      }
    }
    answers {
      answer = "198.51.132.228"
        meta {
          feed = "${nsone_datafeed.exampledc3.id}"
          field = "up"
        }
    }
}

resource "nsone_record" "m" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "m.${var.tld}"
    type = "A"
    answers {
      answer = "198.51.132.28"
    }
    answers {
      answer = "198.51.132.228"
    }
}

resource "nsone_record" "m_star" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "*.m.${var.tld}"
    type = "CNAME"
    answers {
      answer = "m.${var.tld}."
    }
}

resource "nsone_record" "biz" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "biz.${var.tld}"
    type = "A"
    answers {
      answer = "198.51.132.28"
    }
    answers {
      answer = "198.51.132.228"
    }
}

resource "nsone_record" "biz_star" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "*.biz.${var.tld}"
    type = "CNAME"
    answers {
      answer = "biz.${var.tld}."
    }
}

resource "nsone_record" "business" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "business.${var.tld}"
    type = "CNAME"
    answers {
      answer = "biz.${var.tld}."
    }
}
    
resource "nsone_record" "adsp_domainkey" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "_adsp._domainkey.${var.tld}"
    type = "TXT"
    answers {
      answer = "dkim=unknown;"
    }
}

resource "nsone_record" "dmarc" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "_dmarc.${var.tld}"
    type = "TXT"
    answers {
      answer = "v=DMARC1; p=none; rua=mailto:dmarc@yelp.com; ruf=mailto:dmarc@yelp.com; ri=43200;"
    }
}

resource "nsone_record" "domainkey" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "_domainkey.${var.tld}"
    type = "TXT"
    answers {
      answer = "o=~;"
    }
}

resource "nsone_record" "google_domainkey" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "google._domainkey.${var.tld}"
    type = "TXT"
    answers {
      answer = "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCE461lSdBhlo1OjyAZttHGZ9wovs2sUmLVE268UOUoZ6tsFKsaQjySIjmRyKTckO7Pca0BcYdqaLW0qN1S8ELf37Bcz4AmdOb4I5uN2QhY3JypbJfb3PNmPcR0mZeZrE9BAADkEHAv0hUimJMcqMPP0+S6SWeLWIJ+lo4LGBmQKwIDAQAB;"
    }
}

resource "nsone_record" "spf1" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "_spf1.${var.tld}"
    type = "TXT"
    answers {
        answer = "v=spf1 ip4:198.51.132.21 ip4:199.255.189.136/29 ip4:174.36.206.82 ip4:72.18.233.188/30 ip4:199.255.189.5 ip4:199.255.189.180/30 ip4:209.58.196.120 a:admin.yelpcorp.com ip4:209.177.160.9 ip4:209.177.164.34 ip4:209.177.168.2 include:_spf1.yelp.com -all"
    }
}

resource "nsone_record" "spf" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "${var.tld}"
    type = "TXT"
    answers {
      answer = "321815654-6930283"
    }
    answers {
      answer = "v=spf1 ip4:198.51.132.21 ip4:199.255.189.136/29 ip4:174.36.206.82 ip4:72.18.233.188/30 ip4:199.255.189.5 ip4:199.255.189.180/30 ip4:209.58.196.120 a:admin.yelpcorp.com ip4:209.177.160.9 ip4:209.177.164.34 ip4:209.177.168.2 include:_spf1.yelp.com -all"
    }
}

resource "nsone_record" "google" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "googleb6cfac2316ad3b65.${var.tld}"
    type = "CNAME"
    answers {
      answer = "google.com."
    }
}

resource "nsone_record" "mx" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "${var.tld}"
    type = "MX"
    answers {
      answer = "1 ASPMX.L.GOOGLE.COM."
    }
    answers {
      answer = "5 ALT1.ASPMX.L.GOOGLE.COM."
    }
    answers {
      answer = "5 ALT2.ASPMX.L.GOOGLE.COM."
    }
    answers {
      answer = "10 ASPMX2.GOOGLEMAIL.COM."
    }
    answers {
      answer = "10 ASPMX3.GOOGLEMAIL.COM."
    }
    answers {
      answer = "10 ASPMX4.GOOGLEMAIL.COM."
    }
    answers {
      answer = "10 ASPMX5.GOOGLEMAIL.COM."
    }
}

resource "nsone_record" "www_ro" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "www.ro.${var.tld}"
    type = "CNAME"
    answers {
      answer = "www.${var.tld}."
    }
}

resource "nsone_record" "biz_ro" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "biz.ro.${var.tld}"
    type = "CNAME"
    answers {
      answer = "biz.${var.tld}."
    }
}

resource "nsone_record" "m_ro" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "m.ro.${var.tld}"
    type = "CNAME"
    answers {
      answer = "m.${var.tld}."
    }
}

resource "nsone_record" "mobile" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "mobile.${var.tld}"
    type = "CNAME"
    answers {
      answer = "m.${var.tld}."
    }
}

resource "nsone_record" "mkto_domainkey" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "mkto._domainkey.${var.tld}"
    type = "TXT"
    answers {
      answer = "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDlrsUEb2/s7XdpaMlxURe9K0Gj7uRgUiyp8ysIFMxqbLaVVfvWl5KvuYroQFQKPYaORFBkobEko1YH4TeZ4HCrcCl85Y7eimJ6f1xDgStyTSaHOXSbUKVpCEoQNqwxIURnj4XbvQ0FKg7aXL2Bd8fUbWP2nYABaHNZYRFRUXSu5QIDAQAB;"
    }
}

resource "nsone_record" "mail_domainkey" {
    zone = "${nsone_zone.yelp_tld.zone}"
    domain = "mail._domainkey.${var.tld}"
    type = "TXT"
    answers {
      answer = "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChwjAtKSrPDGvMfHExH0dpGwAtYfoG4phx6yZ8IQOR2o5lyDVJRM8px7P/hXy46A1vVvhj1T5+bVnQFdZgzesstEmjERMR+SBXPbGDDKzcoSrICPGxFkTSME5dX953DYz9doqlI/DJ6nd2Bi64kRNHMAM3lNP2Vz+8m7rVB05J/QIDAQAB;"
    }
}

