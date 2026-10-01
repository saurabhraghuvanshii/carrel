#include <string>
using namespace std;

// Recursive descent: read letters and k[...] groups until a ']' or the end.
static string decodeFrom(const string& text, size_t& i) {
    string out;
    while (i < text.size() && text[i] != ']') {
        if (isdigit((unsigned char)text[i])) {
            int k = 0;
            while (isdigit((unsigned char)text[i])) k = k * 10 + (text[i++] - '0');
            i++;  // '['
            string inner = decodeFrom(text, i);
            i++;  // ']'
            for (int r = 0; r < k; r++) out += inner;
        } else {
            out += text[i++];
        }
    }
    return out;
}

string decode(string& text) {
    size_t i = 0;
    return decodeFrom(text, i);
}
