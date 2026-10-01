#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iostream>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

struct ListNode {
    int val;
    ListNode* next;
    ListNode(int v) : val(v), next(nullptr) {}
};

#include "solution.cpp"

// Reads T cases from stdin, each "n" then n values. Prints the returned list
// as [a, b, c] on one line. Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        int n;
        cin >> n;
        ListNode dummy(0);
        ListNode* tail = &dummy;
        for (int i = 0; i < n; i++) {
            int v;
            cin >> v;
            tail->next = new ListNode(v);
            tail = tail->next;
        }
        try {
            ListNode* r = reverseList(dummy.next);
            ostringstream line;
            line << '[';
            int count = 0;
            for (ListNode* p = r; p != nullptr; p = p->next) {
                if (++count > n) throw runtime_error("the returned list has more nodes than the input. Is there a loop?");
                if (count > 1) line << ", ";
                line << p->val;
            }
            out << line.str() << "]\n";
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
