#include <sstream>
#include <string>
#include <vector>
using namespace std;

string reverseWords(string& text) {
    istringstream in(text);
    vector<string> words;
    string w;
    while (in >> w) words.push_back(w);
    string out;
    for (auto it = words.rbegin(); it != words.rend(); ++it) {
        if (!out.empty()) out += ' ';
        out += *it;
    }
    return out;
}
