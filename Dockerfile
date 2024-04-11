FROM httpd:2.4

CMD httpd-foreground -c "LoadModule cgid_module modules/mod_cgid.so"