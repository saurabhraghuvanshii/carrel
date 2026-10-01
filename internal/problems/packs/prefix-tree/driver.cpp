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

// Reads T cases from stdin and prints one line per case. Each case is q, then q operations: "insert w", "search w" or "prefix p".
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        int q;
        cin >> q;
        vector<string> ops(q), arg(q);
        for (int k = 0; k < q; k++) cin >> ops[k] >> arg[k];
        try {
            Trie trie;
            string results = "[";
            for (int k = 0; k < q; k++) {
                if (ops[k] == "insert") {
                    trie.insert(arg[k]);
                    continue;
                }
                bool r = ops[k] == "search" ? trie.search(arg[k]) : trie.startsWith(arg[k]);
                results += (results.size() > 1 ? ", " : "") + string(r ? "true" : "false");
            }
            out << results << "]\n";
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
