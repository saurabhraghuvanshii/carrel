#include <string>
using namespace std;

string squash(string& word) {
    string out;
    for (char c : word)
        if (out.empty() || out.back() != c) out += c;
    return out;
}
