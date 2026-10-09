import importlib.util
import pathlib
import unittest

spec = importlib.util.spec_from_file_location("build_legal", pathlib.Path(__file__).with_name("build-legal.py"))
build_legal = importlib.util.module_from_spec(spec)
spec.loader.exec_module(build_legal)

DETAILS = {"version": "2026-10-09", "service_name": "Shop & Co", "company_name": "Acme <Ltd>", "paid_plans": False, "refund_policy": "No refunds"}


class RenderTests(unittest.TestCase):
    def test_values_are_escaped(self):
        out = build_legal.render_template("Run by {{company_name}} for {{service_name}}", DETAILS)
        self.assertEqual(out, "Run by Acme &lt;Ltd&gt; for Shop &amp; Co")

    def test_placeholders_are_highlighted(self):
        out = build_legal.render_template("{{company_name}}", {**DETAILS, "company_name": "REPLACE: name"})
        self.assertEqual(out, '<mark class="todo">REPLACE: name</mark>')

    def test_plain_values_never_get_a_highlight_inside_an_attribute(self):
        details = {**DETAILS, "contact_email": "REPLACE: email"}
        out = build_legal.render_template('<a href="mailto:{{contact_email|plain}}">{{contact_email}}</a>', details)
        self.assertEqual(out, '<a href="mailto:REPLACE: email"><mark class="todo">REPLACE: email</mark></a>')
        self.assertNotIn("<mark", out.split(">")[0])

    def test_unknown_key_is_an_error(self):
        with self.assertRaises(build_legal.BuildError):
            build_legal.render_template("{{nope}}", DETAILS)

    def test_free_or_paid_wording(self):
        template = "<!--IF paid-->PAID<!--ELSE-->FREE<!--END-->"
        self.assertEqual(build_legal.render_template(template, DETAILS), "FREE")
        self.assertEqual(build_legal.render_template(template, {**DETAILS, "paid_plans": True}), "PAID")

    def test_paid_plans_need_a_refund_policy(self):
        with self.assertRaises(build_legal.BuildError):
            build_legal.validate({**DETAILS, "paid_plans": True, "refund_policy": "REPLACE: later"})
        build_legal.validate({**DETAILS, "paid_plans": True})

    def test_version_must_be_short_and_plain(self):
        for bad in ("", "x" * 33, "has space", "<script>"):
            with self.assertRaises(build_legal.BuildError):
                build_legal.validate({**DETAILS, "version": bad})


class RealTemplateTests(unittest.TestCase):
    def test_no_html_element_ever_appears_inside_an_attribute(self):
        details = build_legal.json.loads((build_legal.LEGAL / "details.json").read_text())
        for key in details:
            if not key.startswith("_") and key != "version" and isinstance(details[key], str):
                details[key] = "REPLACE: " + key  # worst case: everything unfilled
        for path, content in build_legal.build(details).items():
            self.assertNotRegex(content, r'="[^"]*<mark', path.name)

    def test_every_placeholder_in_the_real_templates_has_a_value(self):
        details = build_legal.json.loads((build_legal.LEGAL / "details.json").read_text())
        files = build_legal.build(details)  # raises if a template uses a key that details.json lacks
        self.assertEqual({p.name for p in files}, {"terms.html", "privacy.html", "legalVersion.ts"})
        for path, content in files.items():
            self.assertNotIn("{{", content, path.name)
            self.assertNotIn("<!--IF", content, path.name)


if __name__ == "__main__":
    unittest.main()
