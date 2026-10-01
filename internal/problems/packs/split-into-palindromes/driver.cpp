#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iomanip>
#include <iostream>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <stdexcept>
#include <stack>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

#include "solution.cpp"

string quoted(const vector<string>& words) {
    string out = "[";
    for (size_t i = 0; i < words.size(); i++) out += (i ? ", \"" : "\"") + words[i] + "\"";
    return out + "]";
}

// Sorts the words, so any order is accepted, then prints ["a", "b"].
string showWords(vector<string> words) {
    sort(words.begin(), words.end());
    return quoted(words);
}

// Sorts the lists (not the words inside them), then prints [["a", "b"], ["ab"]].
string showSplits(vector<vector<string>> lists) {
    sort(lists.begin(), lists.end());
    string out = "[";
    for (size_t i = 0; i < lists.size(); i++) out += (i ? ", " : "") + quoted(lists[i]);
    return out + "]";
}

// Reads T cases from stdin and prints one line per case. Each case is one word.
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        string word;
        cin >> word;
        try {
            out << showSplits(cuts(word)) << '\n';
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
